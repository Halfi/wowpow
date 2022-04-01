package server

import (
	"time"
)

type options struct {
	listenersLimit int64
	timeout        time.Duration
	bits           int32
	secret         string
}

type Options func(*options)

func WithListenersLimit(callback int64) Options {
	return func(options *options) {
		options.listenersLimit = callback
	}
}

func WithBits(bits int32) Options {
	return func(options *options) {
		options.bits = bits
	}
}

func WithSecret(secret string) Options {
	return func(options *options) {
		options.secret = secret
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(options *options) {
		options.timeout = timeout
	}
}
