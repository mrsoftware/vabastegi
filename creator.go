package app

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// Creator is what a pkgs use to create dependency.
type Creator[T any] func(ctx context.Context, hub T) error

// CreatorErrorBlocker   used to check if we should stop in error cases or not.
type CreatorErrorBlocker[T any] interface {
	BlockOnError() (block bool)
	Create(ctx context.Context, hub T) error
}

// Create call all passed creators.
// creators must be Creator or CreatorErrorBlocker, otherwise return error on call.
func Create[T any](ctx context.Context, creator ...interface{}) (*Depend[T], error) {
	depend := &Depend[T]{waitGroup: sync.WaitGroup{}}

	for _, create := range creator {
		fmt.Printf("Running %s ", getCreatorName(create))

		err := call(ctx, depend, create)
		if err != nil {
			fmt.Printf("✕\n")
		}

		if err != nil && shouldBlock[T](create) {
			return nil, err
		}

		if err != nil {
			doLog(depend, err)
		}

		fmt.Printf("✓\n")
	}

	return depend, nil
}

func call[T any](ctx context.Context, depend *Depend[T], creator interface{}) error {
	callable, isOk := creator.(Creator[T])
	if isOk {
		return callable(ctx, depend.hub)
	}

	callableMethod, ok := creator.(CreatorErrorBlocker[T])
	if ok {
		return callableMethod.Create(ctx, depend.hub)
	}

	return fmt.Errorf("passed creator is not valid: %#v", creator)
}

func doLog[T any](hub *Depend[T], err error) {
	log.Printf("creator failed due to: %#v", err)
}

func shouldBlock[T any](creator interface{}) bool {
	blocker, ok := creator.(CreatorErrorBlocker[T])
	if !ok {
		return true
	}

	return blocker.BlockOnError()
}

func getCreatorName(creator interface{}) string {
	parts := strings.Split(GetFunctionNameReference(creator), ".")

	return parts[len(parts)-1]
}

// GetFunctionNameReference of passed arg.
func GetFunctionNameReference(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
