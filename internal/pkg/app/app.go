package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var ErrShutdownTimeout = fmt.Errorf("application shutdown forced")

// Application - container of runnable objects with run logic and graceful shutdown
// It has graceful shutdown mechanism. Application wait signal SIGINT.
// When application catches SIGINT signal it cancel main context.
// If after cancel all runners won't be stopped in shutdownTimeout interval, application will be force closed.
type Application struct {
	mu              sync.Mutex
	shutdownTimeout time.Duration

	runnables []Runner
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

// Runner interface which could be run in app
type Runner interface {
	// Run starts runner working
	// Function should be blocking. All functions should catch ctx.Done() channel and correctly finalize their work.
	// If it returns nil, it means that application runner finished work successfully
	// If it returns error, it means that runner has exceptionally problem and application should be stopped
	Run(ctx context.Context) error
}

// Register register run class
func (a *Application) Register(r Runner) {
	a.runnables = append(a.runnables, r)
}

// Run application
// This is blocking function. It will unblock in two keyses:
//   - Any runner stopped with error
//   - All runners stopped
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
	a.mu.Lock()
	defer a.mu.Unlock()

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

func newRunFn(ctx context.Context, r Runner) func() error {
	return func() error {
		return r.Run(ctx)
	}
}
