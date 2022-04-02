package server

//go:generate mockgen -package=mock -destination=./mock/server.go wowpow/internal/pkg/server Verifier,Messenger,Listener,Conn

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/proto"

	"wowpow/internal/pkg/hash"
	"wowpow/internal/pkg/pow"
	"wowpow/pkg/api/message"
)

var (
	ErrConnectionClose = fmt.Errorf("close")
	ErrUnknownRequest  = fmt.Errorf("unknown request")
)

const (
	nl             byte = 10 // ascii code of new line
	defaultTimeout      = time.Minute
)

// Verifier verify hashcash
type Verifier interface {
	Verify(ctx context.Context, h *pow.Hashcach, resource string) error
}

// Messenger returns random quote
type Messenger interface {
	GetMessage() string
}

// Listener is an wrapper to create mocks
type Listener interface {
	net.Listener
}

// Conn is an wrapper to create mocks
type Conn interface {
	net.Conn
}

// Server main tcp server struct
type Server struct {
	lis       Listener
	hasher    hash.Hasher
	verifier  Verifier
	messenger Messenger
	options   options
}

// New constructor
func New(lis Listener, hasher hash.Hasher, verifier Verifier, messenger Messenger, opts ...Options) *Server {
	s := &Server{lis: lis, hasher: hasher, verifier: verifier, messenger: messenger}

	for i := range opts {
		opts[i](&s.options)
	}

	if s.options.listenersLimit == 0 {
		s.options.listenersLimit = int64(runtime.GOMAXPROCS(0))
	}

	if s.options.timeout == 0 {
		s.options.timeout = defaultTimeout
	}

	return s
}

// Run main server loop.
func (s *Server) Run(ctx context.Context) error {
	log.Printf("listen tcp on addr %s", s.lis.Addr())
	payloadChan := make(chan Conn, 1)
	sem := semaphore.NewWeighted(s.options.listenersLimit)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			if err := ctx.Err(); err != nil {
				return
			}

			conn, err := s.lis.Accept()
			if err != nil {
				// connection closed. HTTP2 package detects closed connection same way.
				if strings.Contains(err.Error(), "use of closed network connection") {
					cancel()
					return
				}

				log.Printf("connection accept error: %s", err)
				continue
			}

			payloadChan <- conn
		}
	}()

	for {
		select {
		case conn := <-payloadChan:
			if err := sem.Acquire(ctx, 1); err != nil {
				return fmt.Errorf("failed to acquire semaphore: %w", err)
			}

			if err := conn.SetDeadline(time.Now().Add(s.options.timeout)); err != nil {
				return fmt.Errorf("failed to set server connection deadline: %w", err)
			}

			go func() {
				defer sem.Release(1)
				s.handle(ctx, conn)
			}()
		case <-ctx.Done():
			// TODO implement correct client finalizing. This is poc and test task.
			// Here we can just close listener and all connections will be broken.
			return s.lis.Close()
		}
	}
}

func (s *Server) handle(ctx context.Context, conn Conn) {
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString(nl)
		if err != nil {
			break
		}

		remoteAddr := conn.RemoteAddr()
		addr := remoteAddr.(*net.TCPAddr)
		if addr == nil {
			fmt.Printf("wrong addr type: %+v", remoteAddr)
			break
		}

		msg, err := s.process(ctx, req, addr.IP.String())
		if err != nil {
			if !errors.Is(err, ErrConnectionClose) {
				fmt.Printf("err process request: %s", err)
			}
			break
		}

		if msg != nil {
			err := s.response(msg, conn)
			if err != nil {
				fmt.Printf("err send message: %s", err)
				break
			}
		}
	}
}

func (s *Server) process(ctx context.Context, req, resource string) (*message.Message, error) {
	req = strings.Trim(req, string(nl))
	msg := new(message.Message)

	buf, err := base64.RawStdEncoding.DecodeString(req)
	if err != nil {
		return nil, fmt.Errorf("decode hex message error: %w", err)
	}

	err = proto.Unmarshal(buf, msg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal proto message error: %w", err)
	}

	switch msg.Header {
	case message.Message_close:
		return s.close()
	case message.Message_challenge:
		return s.challenge(msg, resource)
	case message.Message_resource:
		return s.resource(ctx, msg, resource)
	default:
		return nil, ErrUnknownRequest
	}
}

func (s *Server) close() (*message.Message, error) {
	return nil, ErrConnectionClose
}

func (s *Server) challenge(msg *message.Message, resource string) (*message.Message, error) {
	hashcash, err := pow.InitHashcash(s.options.bits, resource, s.options.secret, s.hasher)
	if err != nil {
		return nil, fmt.Errorf("init hashcash error: %w", err)
	}

	m := hashcash.String()
	_ = m
	msg.Response = &message.Message_Hashcach{Hashcach: hashcash.ToProto()}

	return msg, nil
}

func (s *Server) resource(ctx context.Context, msg *message.Message, resource string) (*message.Message, error) {
	resp := msg.GetHashcach()
	if resp != nil {
		hashcash, err := pow.FromProto(resp)
		if err != nil {
			log.Printf("parse hashcash from respone error: %s", err)
			return nil, ErrUnknownRequest
		}

		err = s.verifier.Verify(ctx, hashcash, resource)
		if err != nil {
			log.Printf("parse hashcash from respone error: %s", err)
			return nil, ErrUnknownRequest
		}

		msg.Response = &message.Message_Payload{Payload: s.messenger.GetMessage()}

		return msg, nil
	}

	return nil, ErrUnknownRequest
}

func (s *Server) response(msg *message.Message, writer io.Writer) error {
	bin, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("server proto message response parse error: %w", err)
	}

	buf := make([]byte, base64.RawStdEncoding.EncodedLen(len(bin)))
	base64.RawStdEncoding.Encode(buf, bin)

	_, err = writer.Write(buf)
	if err != nil {
		return fmt.Errorf("server send message response error: %w", err)
	}

	_, err = writer.Write([]byte{nl})
	if err != nil {
		return fmt.Errorf("server send finalyze response message error: %w", err)
	}

	return nil
}
