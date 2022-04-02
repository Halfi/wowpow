package dialer

import (
	"fmt"
	"net"
)

type TCP struct {
	address string
}

// New constructor
func New(address string) *TCP {
	tcp := &TCP{
		address: address,
	}

	return tcp
}

// Dial make tcp dial
func (t *TCP) Dial() (Conn, error) {
	conn, err := net.Dial("tcp", t.address)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	return &connection{Conn: conn}, nil
}
