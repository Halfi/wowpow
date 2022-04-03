package client

import (
	"encoding/base64"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/proto"

	"wowpow/internal/pkg/dialer"
	"wowpow/pkg/api/message"
)

// Connection manager manage all connection through dialer.Conn interface
type Connection struct {
	conn dialer.Conn

	mu       sync.RWMutex
	done     chan interface{}
	send     chan *message.Message
	response *string
	err      error
}

// NewConnection constructor sets default values
func NewConnection(conn dialer.Conn, maxProc int64) *Connection {
	c := &Connection{
		conn: conn,
		done: make(chan interface{}),
		send: make(chan *message.Message),
	}

	sem := semaphore.NewWeighted(maxProc)

	go func() {
		for {
			if c.isClosed() {
				return
			}

			if sem.TryAcquire(1) {
				go func() {
					defer func() {
						sem.Release(1)
					}()

					if err := c.sending(); err != nil {
						// on sending error set error and stop the loop
						c.SetResponse("", err)
					}
				}()
			}
		}
	}()

	return c
}

// GetConnection returns tcp connection
func (c *Connection) GetConnection() dialer.Conn {
	return c.conn
}

// SendInitRequest sends first challenge request
func (c *Connection) SendInitRequest() {
	c.Send(&message.Message{
		Header: message.Message_challenge,
	})
}

// SendCloseMessage sends close message to the server
func (c *Connection) SendCloseMessage() {
	c.Send(&message.Message{
		Header: message.Message_close,
	})
}

// Send message to the server
func (c *Connection) Send(m *message.Message) {
	if !c.isClosed() {
		c.send <- m
	}
}

// SetResponse setting response and close connection and all channels
func (c *Connection) SetResponse(res string, err error) {
	c.SendCloseMessage()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.response == nil {
		c.response = &res
		c.err = err
		close(c.done)
		close(c.send)
		_ = c.conn.Close()
	}
}

// Done channel which will be closed on finish
func (c *Connection) Done() <-chan interface{} {
	return c.done
}

// GetResponse returns response
func (c *Connection) GetResponse() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var res string
	if c.response != nil {
		res = *c.response
	}

	return res, c.err
}

func (c *Connection) sending() error {
	for req := range c.send {
		bin, err := proto.Marshal(req)
		if err != nil {
			return fmt.Errorf("client proto message parse error: %w", err)
		}

		buf := make([]byte, base64.RawStdEncoding.EncodedLen(len(bin)))
		base64.RawStdEncoding.Encode(buf, bin)

		_, err = c.conn.Write(buf)
		if err != nil {
			return fmt.Errorf("client send message error: %w", err)
		}

		// send new line to finalize request
		_, err = c.conn.Write([]byte{nl})
		if err != nil {
			return fmt.Errorf("client send message error: %w", err)
		}
	}

	return nil
}

func (c *Connection) isClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.response != nil
}
