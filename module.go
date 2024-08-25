package vabastegi

import "context"

func NewModule[T any](provider ModuleProvider[T]) Provider[T] {
	return func(ctx context.Context, a *App[T]) error {
		module := provider(ctx, a)
		if module.Err != nil {
			return module.Err
		}

		return SetHubField(a, module.Name, module.Data)
	}
}

type ModuleProvider[T any] func(ctx context.Context, app *App[T]) Module

type Module struct {
	Name string
	Data interface{}
	Err  error
}
