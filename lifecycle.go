package vabastegi

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/mrsoftware/errors"
)

// ShutdownListener is what you need to pass to OnShutdown method.
type ShutdownListener = func(ctx context.Context) error

// lifecycle manage application lifecycle like running task or shutdown.
type lifecycle struct {
	listeners []ShutdownListener
	publisher *eventManager
	ctx       context.Context
	cancel    context.CancelCauseFunc
	waitGroup *errors.WaitGroup
}

// newLifecycle create a new instance of lifecycle.
func newLifecycle(ctx context.Context, publisher *eventManager) *lifecycle {
	waitGroup := errors.NewWaitGroup(errors.WaitGroupWithContext(ctx), errors.WaitGroupWithStopOnError())

	return &lifecycle{
		listeners: make([]ShutdownListener, 0),
		publisher: publisher,
		ctx:       waitGroup.Context(),
		cancel:    waitGroup.Stop,
		waitGroup: waitGroup,
	}
}

// Wait on shutdown or application finish.
func (l *lifecycle) Wait() error {
	var err error
	select {
	case <-l.ctx.Done():
		err = l.ctx.Err()
	case err = <-errors.WaitChanel(l.waitGroup):
		l.Stop(err)
	}

	return l.callShutdownListeners(l.ctx, err)
}

// RegisterGracefulShutdown start listing on os signal and cancel the parent context on getting one.
func (l *lifecycle) RegisterGracefulShutdown() {
	ctx, cancel := signal.NotifyContext(l.ctx, os.Interrupt)

	l.ctx, l.cancel = ctx, cancelToCancelCause(cancel)
}

// Stop the application.
func (l *lifecycle) Stop(cause error) {
	l.cancel(cause)
}

// GetContext of lifecycle.
func (l *lifecycle) GetContext() context.Context {
	return l.ctx
}

// do Shut down the application.
func (l *lifecycle) callShutdownListeners(ctx context.Context, cause error) error {
	errList := errors.NewMultiError()

	startAt := time.Now()

	l.publisher.Publish(&OnApplicationShutdownExecuting{
		Cause:      cause,
		ShutdownAt: startAt,
	})

	defer func() {
		l.publisher.Publish(&OnApplicationShutdownExecuted{
			Runtime: time.Now().Sub(startAt),
			Err:     errList.Err(),
		})
	}()

	for _, fn := range l.listeners {
		errList.Add(l.shutdown(ctx, fn))
	}

	return errList.Err()
}

func (l *lifecycle) shutdown(ctx context.Context, callback ShutdownListener) (err error) {
	startAt := time.Now()

	l.publisher.Publish(&OnShutdownExecuting{
		ProviderName: getProviderName(callback, 1),
		CallerPath:   getProviderName(callback, -1),
		ShutdownAt:   startAt,
	})

	defer func() {
		l.publisher.Publish(&OnShutdownExecuted{
			ProviderName: getProviderName(callback, 1),
			CallerPath:   getProviderName(callback, -1),
			Runtime:      time.Now().Sub(startAt),
			Err:          err,
		})
	}()

	return callback(ctx)
}

// OnShutdown add callback to a list of listeners.
func (l *lifecycle) OnShutdown(callback ShutdownListener) {
	l.listeners = append(l.listeners, callback)
}

// RunTask in the background.
func (l *lifecycle) RunTask(ctx context.Context, fn func(ctx context.Context) error) {
	l.waitGroup.DoWithContext(ctx, fn)
}

func getProviderName(creator interface{}, index int) string {
	reference := runtime.FuncForPC(reflect.ValueOf(creator).Pointer()).Name()
	if index == -1 {
		return reference
	}

	parts := strings.Split(reference, ".")

	return parts[len(parts)-(1+index)]
}

// cancelToCancelCause is just a wrapper to turn context.CancelFunc into context.CancelCauseFunc.
func cancelToCancelCause(cancelFunc context.CancelFunc) context.CancelCauseFunc {
	return func(cause error) { cancelFunc() }
}
