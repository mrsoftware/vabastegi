package vabastegi

import (
	"context"
	"io"
	"time"

	"github.com/mrsoftware/errors"
)

// ErrApplicationTimeout .
var ErrApplicationTimeout = errors.New("application timeout")

// Options of Vabastegi.
type Options struct {
	GracefulShutdown bool
	AppName          string
	EventHandlers    []EventHandler
	Ctx              context.Context
}

// Option of Vabastegi.
type Option func(options *Options)

// WithGraceFullShutdown used if you need gracefully shutdown for your application.
func WithGraceFullShutdown(active bool) Option {
	return func(options *Options) {
		options.GracefulShutdown = active
	}
}

// WithAppName provide appName for
func WithAppName(appName string) Option {
	return func(options *Options) {
		options.AppName = appName
	}
}

// WithEventHandlers register event handlers for vabastegi events.
func WithEventHandlers(handlers ...EventHandler) Option {
	return func(options *Options) {
		options.EventHandlers = append(options.EventHandlers, handlers...)
	}
}

// WithContext used if you need to pass custom context.
func WithContext(ctx context.Context) Option {
	return func(options *Options) {
		options.Ctx = ctx
	}
}

// WithTimeout used if you need to pass timout for application booting.
func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		// ignoring cancel
		ctx, _ := context.WithTimeoutCause(options.Ctx, timeout, ErrApplicationTimeout)

		options.Ctx = ctx
	}
}

// WithEventLogger register event handlers for logger.
func WithEventLogger(writer io.Writer, level logLevel) Option {
	return func(options *Options) {
		options.EventHandlers = append(options.EventHandlers, NewEventLogger(writer, level))
	}
}
