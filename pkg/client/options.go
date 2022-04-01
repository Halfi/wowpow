package client

import (
	"time"
)

type options struct {
	timeout       time.Duration
	maxIterations int64
}

type Options func(*options)

func WithTimeout(timeout time.Duration) Options {
	return func(options *options) {
		options.timeout = timeout
	}
}

func WithMaxIterations(maxIterations int64) Options {
	return func(options *options) {
		options.maxIterations = maxIterations
	}
}
