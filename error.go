package depend

import (
	"context"
	"fmt"
)

// IgnoreError is a provider wrapper and used if you want to ignore provider error.
func IgnoreError[T any](provider Provider[T]) Provider[T] {
	return func(ctx context.Context, depend *App[T]) error {
		if err := provider(ctx, depend); err != nil {
			depend.Logger().Error(fmt.Sprintf("privider error ignored: %s", err.Error()))
		}

		return nil
	}
}
