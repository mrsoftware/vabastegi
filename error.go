package vabastegi

import (
	"context"
)

// IgnoreError is a provider wrapper and used if you want to ignore provider error.
func IgnoreError[T any](provider Provider[T]) Provider[T] {
	return func(ctx context.Context, app *App[T]) error {
		if err := provider(ctx, app); err != nil {
			app.Log(ErrorLogLevel, "provider error ignored: %s", err)
		}

		return nil
	}
}
