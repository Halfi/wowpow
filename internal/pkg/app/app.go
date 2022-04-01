package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

var ErrShutdownTimeout = fmt.Errorf("application shutdown forced")

// Application - container of runnable objects with run logic and graceful shutdown
type Application struct {
	shutdownTimeout time.Duration

	runnables []Runnable
	cancel    context.CancelFunc
	force     chan interface{}
}

// New constructor
func New(shutdownTimeout time.Duration) *Application {
	return &Application{
		shutdownTimeout: shutdownTimeout,
		force:           make(chan interface{}, 1),
	}
}

type Runnable interface {
	Run(ctx context.Context) error
}

func (a *Application) Register(r Runnable) {
	a.runnables = append(a.runnables, r)
}

func (a *Application) Run(ctx context.Context) error {
	ctx, a.cancel = context.WithCancel(ctx)

	eg, ctx := errgroup.WithContext(ctx)

	go func() {
		<-ctx.Done()
		a.stop()
	}()

	for _, r := range a.runnables {
		eg.Go(newRunFn(ctx, r))
	}

	go a.waitForInterruption()

	return a.wait(eg)
}

func (a *Application) wait(eg *errgroup.Group) error {
	egErr := make(chan error)
	go func() {
		egErr <- eg.Wait()
	}()

	for {
		select {
		case err := <-egErr:
			return err
		case <-a.force:
			return ErrShutdownTimeout
		}
	}
}

func (a *Application) stop() {
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
		go func() {
			timer := time.NewTimer(a.shutdownTimeout)
			defer timer.Stop()

			<-timer.C
			close(a.force)
		}()
	}
}

func (a *Application) waitForInterruption() {
	c := make(chan os.Signal, 1)
	defer func() {
		signal.Stop(c)
		close(c)
	}()

	signal.Notify(c, os.Interrupt)
	<-c
	a.stop()
}

func newRunFn(ctx context.Context, r Runnable) func() error {
	return func() error {
		return r.Run(ctx)
	}
}
