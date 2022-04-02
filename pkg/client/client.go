package client

//go:generate mockgen -package=mock -destination=./mock/client.go wowpow/pkg/client Dialer,Computer

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	"wowpow/internal/pkg/dialer"
	"wowpow/internal/pkg/pow"
	"wowpow/pkg/api/message"
)

const (
	defaultTimeout            = time.Minute
	defaultMaxIterations      = 1 << 20
	nl                   byte = 10 // ascii code of new line
)

// Dialer sends calls throught any protocol (tcp)
type Dialer interface {
	Dial() (dialer.Conn, error)
}

// Computer waste time
type Computer interface {
	Compute(context.Context, *pow.Hashcach, int64) (*pow.Hashcach, error)
}

var (
	ErrConnectionClose = fmt.Errorf("close")
	ErrUnknownResponse = fmt.Errorf("unknown response")
)

type WoWPoW struct {
	computer Computer
	dialer   Dialer
	options  options
}

type connection struct {
	conn       dialer.Conn
	writeCh    chan *message.Message
	responseCh chan string
	errCh      chan error
}

func (conn *connection) initRequest() {
	conn.writeCh <- &message.Message{
		Header: message.Message_challenge,
	}
}

func (conn *connection) close() {
	conn.writeCh <- &message.Message{
		Header: message.Message_close,
	}
}

// NewWoWPoW constructor creates WoWPoW client
func NewWoWPoW(computer Computer, dialer Dialer, opts ...Options) (*WoWPoW, error) {
	w := &WoWPoW{
		computer: computer,
		dialer:   dialer,
	}

	for i := range opts {
		opts[i](&w.options)
	}

	if w.options.timeout == 0 {
		w.options.timeout = defaultTimeout
	}

	if w.options.maxIterations == 0 {
		w.options.maxIterations = defaultMaxIterations
	}

	return w, nil
}

// GetMessage getting quote from remote server with pow challenge
func (w *WoWPoW) GetMessage(ctx context.Context) (string, error) {
	conn, err := w.dialer.Dial()
	if err != nil {
		return "", fmt.Errorf("connection error: %w", err)
	}
	defer func() { _ = conn.Close() }()

	err = conn.SetDeadline(time.Now().Add(w.options.timeout))
	if err != nil {
		return "", fmt.Errorf("connection setting deadline error: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, w.options.timeout)
	defer cancel()

	connection := &connection{
		conn:       conn,
		writeCh:    make(chan *message.Message, 1),
		responseCh: make(chan string, 1),
		errCh:      make(chan error, 1),
	}

	connection.initRequest()
	defer connection.close()

	go func(ctx context.Context) {
		<-ctx.Done()
		connection.errCh <- ctx.Err()
	}(ctx)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return w.read(ctx, connection)
	})

	eg.Go(func() error {
		return w.write(ctx, connection)
	})

	go func() {
		err := eg.Wait()
		if err != nil {
			connection.errCh <- err
		}
	}()

	select {
	case response := <-connection.responseCh:
		return response, nil
	case err := <-connection.errCh:
		return "", err
	}
}

func (w *WoWPoW) read(ctx context.Context, conn *connection) error {
	defer func() { _ = conn.conn.Close() }()

	reader := bufio.NewReader(conn.conn)
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("reading error: %w", err)
		}

		res, err := reader.ReadString(nl)
		if err != nil {
			return fmt.Errorf("reading error: %w", err)
		}

		err = w.process(ctx, res, conn)
		if err != nil {
			return fmt.Errorf("client read process error: %w", err)
		}
	}
}

func (w *WoWPoW) process(ctx context.Context, res string, conn *connection) error {
	res = strings.Trim(res, string(nl))
	msg := new(message.Message)

	buf, err := base64.RawStdEncoding.DecodeString(res)
	if err != nil {
		return fmt.Errorf("decode hex message error: %w", err)
	}

	err = proto.Unmarshal(buf, msg)
	if err != nil {
		return fmt.Errorf("unmarshal response proto message error: %w", err)
	}

	switch msg.Header {
	case message.Message_close:
		return w.close()
	case message.Message_challenge:
		return w.challenge(ctx, msg, conn)
	case message.Message_resource:
		conn.responseCh <- msg.GetPayload()
		return nil
	default:
		return ErrUnknownResponse
	}
}

func (w *WoWPoW) close() error {
	return ErrConnectionClose
}

func (w *WoWPoW) challenge(ctx context.Context, msg *message.Message, conn *connection) error {
	hashcash := msg.GetHashcach()
	if hashcash == nil {
		return ErrUnknownResponse
	}

	hc, err := pow.FromProto(hashcash)
	if err != nil {
		return fmt.Errorf("unmarshal response proto message error: %w", err)
	}

	hc, err = w.computer.Compute(ctx, hc, w.options.maxIterations)
	if err != nil {
		return fmt.Errorf("compute challenge error: %w", err)
	}

	conn.writeCh <- &message.Message{
		Header:   message.Message_resource,
		Response: &message.Message_Hashcach{Hashcach: hc.ToProto()},
	}

	return nil
}

func (w *WoWPoW) write(ctx context.Context, connection *connection) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-connection.writeCh:
			bin, err := proto.Marshal(req)
			if err != nil {
				return fmt.Errorf("client proto message parse error: %w", err)
			}

			buf := make([]byte, base64.RawStdEncoding.EncodedLen(len(bin)))
			base64.RawStdEncoding.Encode(buf, bin)

			_, err = connection.conn.Write(buf)
			if err != nil {
				return fmt.Errorf("client send message error: %w", err)
			}

			// send new line to finalize request
			_, err = connection.conn.Write([]byte{nl})
			if err != nil {
				return fmt.Errorf("client send message error: %w", err)
			}
		}
	}
}
