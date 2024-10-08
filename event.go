package vabastegi

import "time"

// Event is what Vabastegi event looks like.
type Event interface {
	event() // it's private to prevent outside implementation.
}

// eventManager responsible to manage the event system.
type eventManager struct {
	handlers []EventHandler
}

// newEventManager create a new instance of eventManager.
func newEventManager(handlers []EventHandler) *eventManager {
	return &eventManager{handlers: handlers}
}

// Publish passed event using event handlers.
func (e *eventManager) Publish(event Event) {
	for _, handler := range e.handlers {
		handler.OnEvent(event)
	}
}

// Register event handler.
func (e *eventManager) Register(handler EventHandler) {
	e.handlers = append(e.handlers, handler)
}

// EventHandler used if you need to handle the events.
type EventHandler interface {
	OnEvent(event Event)
}

func (p *OnBuildsExecuting) event()              {}
func (p *OnBuildsExecuted) event()               {}
func (p *OnBuildExecuting) event()               {}
func (p *OnBuildExecuted) event()                {}
func (p *OnShutdownExecuting) event()            {}
func (p *OnShutdownExecuted) event()             {}
func (p *OnApplicationShutdownExecuting) event() {}
func (p *OnApplicationShutdownExecuted) event()  {}
func (p *OnLog) event()                          {}

// OnBuildsExecuting is emitted before a Builds is executed.
type OnBuildsExecuting struct {
	// BuildAt is the time build happened.
	BuildAt time.Time
}

// OnBuildsExecuted is emitted after a Builds has been executed.
type OnBuildsExecuted struct {
	// Runtime specifies how long it took to run this hook.
	Runtime time.Duration

	// Err is non-nil if the hook failed to execute.
	Err error
}

// OnBuildExecuting is emitted before a Build is executed.
type OnBuildExecuting struct {
	// BuildAt is the time build happened.
	BuildAt time.Time

	// ProviderName is the name of the function that will be executed.
	ProviderName string

	// CallerPath is the path of provider if from.
	CallerPath string
}

// OnBuildExecuted is emitted after a Provider has been executed.
type OnBuildExecuted struct {
	// ProviderName is the name of the function that was executed.
	ProviderName string

	// CallerPath is the path of provider if from.
	CallerPath string

	// Runtime specifies how long it took to run this hook.
	Runtime time.Duration

	// Err is non-nil if the hook failed to execute.
	Err error
}

// OnShutdownExecuting is emitted before a Shutdown is executed.
type OnShutdownExecuting struct {
	// ShutdownAt is the time shutdown happened.
	ShutdownAt time.Time

	// ProviderName is the name of the function that will be executed.
	ProviderName string

	// CallerPath is the path of provider if from.
	CallerPath string
}

// OnShutdownExecuted is emitted after a Shutdown has been executed.
type OnShutdownExecuted struct {
	// ProviderName is the name of the function that was executed.
	ProviderName string

	// CallerPath is the path of provider if from.
	CallerPath string

	// Runtime specifies how long it took to run this hook.
	Runtime time.Duration

	// Err is non-nil if the hook failed to execute.
	Err error
}

// OnApplicationShutdownExecuting is emitted before the application Shutdown is executed.
type OnApplicationShutdownExecuting struct {
	// ShutdownAt is the time shutdown happened.
	ShutdownAt time.Time

	// Cause is the reason for shutdown the application.
	Cause error
}

// OnApplicationShutdownExecuted is emitted after the application Shutdown has been executed.
type OnApplicationShutdownExecuted struct {
	// Runtime specifies how long it took to run this hook.
	Runtime time.Duration

	// Err is non-nil if the hook failed to execute.
	Err error
}

// OnLog is used if a log event is sent.
type OnLog struct {
	LogAt   time.Time
	Level   logLevel
	Message string
	Args    []interface{}
}
