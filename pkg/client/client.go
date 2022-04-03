package client

//go:generate mockgen -package=mock -destination=./mock/client.go wowpow/pkg/client Dialer,Computer

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"

	"wowpow/internal/pkg/dialer"
	"wowpow/internal/pkg/pow"
	"wowpow/pkg/api/message"
)

const (
	nl byte = 10 // ascii code of new line
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

// NewWoWPoW constructor creates WoWPoW client
func NewWoWPoW(computer Computer, dialer Dialer, opts ...Options) (*WoWPoW, error) {
	w := &WoWPoW{
		computer: computer,
		dialer:   dialer,
	}

	for i := range opts {
		opts[i](&w.options)
	}

	InitDefaultOptions(&w.options)

	return w, nil
}

// GetMessage getting quote from remote server with pow challenge
func (w *WoWPoW) GetMessage(ctx context.Context) (string, error) {
	netConn, err := w.dialer.Dial()
	if err != nil {
		return "", fmt.Errorf("connection error: %w", err)
	}
	defer func() { _ = netConn.Close() }()

	err = netConn.SetDeadline(time.Now().Add(w.options.timeout))
	if err != nil {
		return "", fmt.Errorf("connection setting deadline error: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, w.options.timeout)
	defer cancel()

	conn := NewConnection(netConn, w.options.maxProc)

	conn.SendInitRequest()

	go func(ctx context.Context) {
		<-ctx.Done()
		conn.SetResponse("", ctx.Err())
	}(ctx)

	go func() {
		err := w.read(ctx, conn)
		if err != nil {
			conn.SetResponse("", err)
		}
	}()

	<-conn.Done()
	return conn.GetResponse()
}

func (w *WoWPoW) read(ctx context.Context, conn *Connection) error {
	reader := bufio.NewReader(conn.GetConnection())
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

func (w *WoWPoW) process(ctx context.Context, res string, conn *Connection) error {
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
		conn.SetResponse(msg.GetPayload(), nil)
		return nil
	default:
		return ErrUnknownResponse
	}
}

func (w *WoWPoW) close() error {
	return ErrConnectionClose
}

func (w *WoWPoW) challenge(ctx context.Context, msg *message.Message, conn *Connection) error {
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

	conn.Send(&message.Message{
		Header:   message.Message_resource,
		Response: &message.Message_Hashcach{Hashcach: hc.ToProto()},
	})

	return nil
}
