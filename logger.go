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

var _ Logger = &IOLogger{}

// Logger is what Vabastegi use as logger.
type Logger interface {
	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Errorf(message string, args ...interface{})
}

// NewIOLogger create new instance of IOLogger.
func NewIOLogger(writer io.Writer, level logLevel) Logger {
	return &IOLogger{writer: writer, level: level}
}

// IOLogger is a Logger that write to passed io.Writer.
type IOLogger struct {
	writer io.Writer
	level  logLevel
}

func (i *IOLogger) IsEnable(level logLevel) bool {
	return level >= i.level
}

func (i *IOLogger) Debugf(message string, args ...interface{}) {
	i.log(DebugLogLevel, message, args...)
}

func (i *IOLogger) Infof(message string, args ...interface{}) {
	i.log(InfoLogLevel, message, args...)
}

func (i *IOLogger) Errorf(message string, args ...interface{}) {
	i.log(ErrorLogLevel, message, args...)
}

func (i *IOLogger) log(level logLevel, message string, args ...interface{}) {
	if !i.IsEnable(level) {
		return
	}

	fmt.Fprintf(i.writer, "[Vabastegi] ["+level.String()+"] "+message+"\n", args...)
}
