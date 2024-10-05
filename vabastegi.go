package vabastegi

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mrsoftware/errors"
)

// Provider is a dependency provider for application.
type Provider[T any] func(context.Context, *App[T]) error

// App is the dependency injection manger.
type App[T any] struct {
	waitGroup            sync.WaitGroup
	errors               *errors.MultiError
	onShutdown           []func(ctx context.Context) error
	Hub                  T
	options              Options
	backgroundTasksCount int
	graceFullOnce        sync.Once
	Event                *EventManager
}

// New instance of App Dependency management.
func New[T any](options ...Option) *App[T] {
	var op Options

	for _, option := range options {
		option(&op)
	}

	app := App[T]{
		options: op,
		errors:  errors.NewMultiError(),
		Event:   NewEventManager(op.EventHandlers),
	}

	if op.GracefulShutdown {
		app.registerGracefulShutdown()
	}

	return &app
}

// Builds the dependency structure of your app.
func (a *App[T]) Builds(ctx context.Context, providers ...Provider[T]) (err error) {
	startAt := time.Now()

	a.Event.Publish(&OnBuildsExecuting{BuildAt: startAt})

	defer func() {
		a.Event.Publish(&OnApplicationShutdownExecuted{Runtime: time.Since(startAt), Err: err})
	}()

	for _, provider := range providers {
		err := a.Build(ctx, provider)
		if err == nil {
			continue
		}

		a.Shutdown(ctx, "Provider Failure")

		return err
	}

	return nil
}

// Build use the provider to set a dependency.
func (a *App[T]) Build(ctx context.Context, provider Provider[T]) (err error) {
	startAt := time.Now()

	a.Event.Publish(&OnBuildExecuting{
		ProviderName: a.getProviderName(provider, 0),
		CallerPath:   a.getProviderName(provider, -1),
		BuildAt:      startAt,
	})

	defer func() {
		a.Event.Publish(&OnBuildExecuted{
			ProviderName: a.getProviderName(provider, 0),
			CallerPath:   a.getProviderName(provider, -1),
			Runtime:      time.Now().Sub(startAt),
			Err:          err,
		})
	}()

	return provider(ctx, a)
}

// RunTask in background.
func (a *App[ـ]) RunTask(fn func()) {
	go fn()

	a.backgroundTasksCount++

	a.waitGroup.Add(1)
}

// Wait for background task to done or any shutdown signal.
func (a *App[ـ]) Wait() error {
	a.waitGroup.Wait()

	return a.errors.Err()
}

// Shutdown ths application.
func (a *App[ـ]) Shutdown(ctx context.Context, reason string) {
	startAt := time.Now()

	a.Event.Publish(&OnApplicationShutdownExecuting{
		Reason:     reason,
		ShutdownAt: startAt,
	})

	defer func() {
		a.Event.Publish(&OnApplicationShutdownExecuted{
			Reason:  reason,
			Runtime: time.Now().Sub(startAt),
			Err:     a.errors.Err(),
		})
	}()

	for _, fn := range a.onShutdown {
		a.errors.Add(a.shutdown(ctx, fn))
	}

	for i := 0; i < a.backgroundTasksCount; i++ {
		a.waitGroup.Done()
	}
}

func (a *App[ـ]) shutdown(ctx context.Context, fn func(context.Context) error) (err error) {
	startAt := time.Now()

	a.Event.Publish(&OnShutdownExecuting{
		ProviderName: a.getProviderName(fn, 1),
		CallerPath:   a.getProviderName(fn, -1),
		ShutdownAt:   startAt,
	})

	defer func() {
		a.Event.Publish(&OnShutdownExecuted{
			ProviderName: a.getProviderName(fn, 1),
			CallerPath:   a.getProviderName(fn, -1),
			Runtime:      time.Now().Sub(startAt),
			Err:          err,
		})
	}()

	return fn(ctx)
}

// OnShutdown register any method for Shutdown method.
func (a *App[ـ]) OnShutdown(fn func(ctx context.Context) error) {
	a.onShutdown = append(a.onShutdown, fn)
}

func (a *App[ـ]) getProviderName(creator interface{}, index int) string {
	reference := runtime.FuncForPC(reflect.ValueOf(creator).Pointer()).Name()
	if index == -1 {
		return reference
	}

	parts := strings.Split(reference, ".")

	return parts[len(parts)-(1+index)]
}

// Log the message.
func (a *App[ـ]) Log(level logLevel, message string, args ...interface{}) {
	a.Event.Publish(&OnLog{LogAt: time.Now(), Level: level, Message: message, Args: args})
}

func (a *App[ـ]) registerGracefulShutdown() {
	a.graceFullOnce.Do(func() {
		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, os.Interrupt)

		go func() {
			appSignal := <-interruptChan
			a.Shutdown(context.Background(), appSignal.String())
		}()
	})
}
