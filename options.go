package vabastegi

// Options of Vabastegi.
type Options struct {
	GracefulShutdown bool
	AppName          string
	EventHandlers    []EventHandler
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
