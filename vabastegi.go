package vabastegi

import (
	"context"
	"time"
)

// Provider is a dependency provider for application.
type Provider[T any] func(context.Context, *App[T]) error

// App is the dependency injection manger.
type App[T any] struct {
	Hub       T
	options   Options
	Events    *eventManager
	Lifecycle *lifecycle
}

// New instance of App Dependency management.
func New[T any](options ...Option) *App[T] {
	op := Options{Ctx: context.Background()}

	for _, option := range options {
		option(&op)
	}

	eManager := newEventManager(op.EventHandlers)
	life := newLifecycle(op.Ctx, eManager)

	if op.GracefulShutdown {
		life.RegisterGracefulShutdown()
	}

	return &App[T]{options: op, Events: eManager, Lifecycle: life}
}

// Builds the dependency structure of your app.
func (a *App[T]) Builds(providers ...Provider[T]) (err error) {
	defer func() {
		if err == nil {
			return
		}

		a.Lifecycle.Stop(err)
	}()

	startAt := time.Now()

	a.Events.Publish(&OnBuildsExecuting{BuildAt: startAt})

	defer func() { a.Events.Publish(&OnApplicationShutdownExecuted{Runtime: time.Since(startAt), Err: err}) }()

	for _, provider := range providers {
		if err = a.Build(a.Lifecycle.GetContext(), provider); err == nil {
			continue
		}

		return err
	}

	return nil
}

// Build use the provider to set a dependency.
func (a *App[T]) Build(ctx context.Context, provider Provider[T]) (err error) {
	startAt := time.Now()

	a.Events.Publish(&OnBuildExecuting{
		ProviderName: getProviderName(provider, 0),
		CallerPath:   getProviderName(provider, -1),
		BuildAt:      startAt,
	})

	defer func() {
		a.Events.Publish(&OnBuildExecuted{
			ProviderName: getProviderName(provider, 0),
			CallerPath:   getProviderName(provider, -1),
			Runtime:      time.Now().Sub(startAt),
			Err:          err,
		})
	}()

	return provider(ctx, a)
}

// Log the message.
func (a *App[ـ]) Log(level logLevel, message string, args ...interface{}) {
	a.Events.Publish(&OnLog{LogAt: time.Now(), Level: level, Message: message, Args: args})
}

// Wait for the application lifecycle to finish.
func (a *App[_]) Wait() error {
	return a.Lifecycle.Wait()
}
