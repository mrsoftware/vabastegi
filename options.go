package vabastegi

// Options of Vabastegi.
type Options struct {
	Logger           Logger
	GracefulShutdown bool
	AppName          string
}

// Option of Vabastegi.
type Option func(options *Options)

// WithLogger provide logger for Vabastegi.
func WithLogger(logger Logger) Option {
	return func(options *Options) {
		options.Logger = logger
	}
}

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
