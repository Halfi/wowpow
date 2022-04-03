package client

import (
	"runtime"
	"time"
)

const (
	defaultTimeout       = time.Minute
	defaultMaxIterations = 1 << 20
)

var defaultMaxProc = int64(runtime.GOMAXPROCS(0))

type options struct {
	timeout       time.Duration
	maxIterations int64
	maxProc       int64
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

func WithMaxProc(maxProc int64) Options {
	return func(options *options) {
		options.maxProc = maxProc
	}
}

func InitDefaultOptions(options *options) {
	if options.maxProc == 0 {
		options.maxProc = defaultMaxProc
	}

	if options.timeout == 0 {
		options.timeout = defaultTimeout
	}

	if options.maxIterations == 0 {
		options.maxIterations = defaultMaxIterations
	}
}
