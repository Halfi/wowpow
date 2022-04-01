package app

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

type RunnerFunc func(ctx context.Context) error

func (f RunnerFunc) Run(ctx context.Context) error {
	return f(ctx)
}

func TestAppSignal(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	app := New(10 * time.Millisecond)

	app.Register(RunnerFunc(func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}))

	var eg errgroup.Group
	eg.Go(func() error {
		return app.Run(ctx)
	})

	<-time.NewTimer(time.Millisecond).C
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	a.Nil(err)

	err = eg.Wait()
	a.ErrorIs(err, context.Canceled)
}

func TestAppShutdownTimeout(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()
	app := New(20 * time.Millisecond)

	app.Register(RunnerFunc(func(ctx context.Context) error {
		select {}
	}))

	var eg errgroup.Group
	eg.Go(func() error {
		return app.Run(ctx)
	})

	<-time.NewTimer(time.Millisecond).C
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	a.Nil(err)

	err = eg.Wait()
	a.ErrorIs(err, ErrShutdownTimeout)
}

func TestAppCTXTimeout(t *testing.T) {
	a := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	app := New(20 * time.Millisecond)

	app.Register(RunnerFunc(func(ctx context.Context) error {
		select {}
	}))

	var eg errgroup.Group
	eg.Go(func() error {
		return app.Run(ctx)
	})

	err := eg.Wait()
	a.ErrorIs(err, ErrShutdownTimeout)
}

func TestAppGoodRunnerCTXTimeout(t *testing.T) {
	a := assert.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	app := New(20 * time.Millisecond)

	app.Register(RunnerFunc(func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}))

	var eg errgroup.Group
	eg.Go(func() error {
		return app.Run(ctx)
	})

	err := eg.Wait()
	a.ErrorIs(err, context.DeadlineExceeded)
}
