package dialer

//go:generate mockgen -package=mock -destination=./mock/dialer.go wowpow/internal/pkg/dialer Conn

import (
	"net"
)

type Conn interface {
	net.Conn
}

type connection struct {
	Conn
}
