package vabastegi

import (
	"fmt"
	"io"
)

const (
	debugLevel = "Debug"
	infoLevel  = "Info"
	errorLevel = "Error"
	warnLevel  = "Warn"
)

type logLevel int

// String version of LogLevel.
func (l logLevel) String() string {
	switch l {
	case DebugLogLevel:
		return debugLevel
	case InfoLogLevel:
		return infoLevel
	case ErrorLogLevel:
		return errorLevel
	case WarnLogLevel:
		return warnLevel
	}

	return "Unknown"
}

const (
	DebugLogLevel logLevel = iota - 1
	InfoLogLevel
	WarnLogLevel
	ErrorLogLevel
)

var _ EventHandler = &EventLogger{}

// EventLogger is use event of vabastegi and log them.
type EventLogger struct {
	writer io.Writer
	level  logLevel
}

// NewEventLogger create new instance of EventLogger.
func NewEventLogger(writer io.Writer, level logLevel) *EventLogger {
	return &EventLogger{writer: writer, level: level}
}

func (l *EventLogger) IsEnable(level logLevel) bool {
	return level >= l.level
}

func (l *EventLogger) log(level logLevel, message string, args ...interface{}) {
	if !l.IsEnable(level) {
		return
	}

	fmt.Fprintf(l.writer, "[Vabastegi] ["+level.String()+"] "+message+"\n", args...)
}

func (l *EventLogger) OnEvent(event Event) {
	switch e := event.(type) {
	case *OnBuildExecuted:
		if e.Err != nil {
			l.log(ErrorLogLevel, e.ProviderName+" ✕")

			return
		}

		l.log(InfoLogLevel, e.ProviderName+" ✓")
	case *OnShutdownExecuting:
		l.log(InfoLogLevel, "Shutting Down %s", e.ProviderName)
	case *OnApplicationShutdownExecuting:
		l.log(InfoLogLevel, "Shutting Down Application: %s", e.Reason)
	case *OnLog:
		l.log(e.Level, e.Message, e.Args...)
	}
}
