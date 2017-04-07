package log

import (
	"context"
	"fmt"
	"runtime"
)

var (
	_logger *logger
)

func init() {
	_logger = new()
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func RegisterHandler(handler Handler, levels ...Level) {
	_logger.RegisterHandler(handler, levels...)
}

// Debug level formatted message.
func Debug(v ...interface{}) {
	e := _logger.newEntry(DebugLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	e := _logger.newEntry(DebugLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Info level formatted message.
func Info(v ...interface{}) {
	e := _logger.newEntry(InfoLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	e := _logger.newEntry(InfoLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Warn level formatted message.
func Warn(v ...interface{}) {
	e := _logger.newEntry(WarnLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	e := _logger.newEntry(WarnLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Error level formatted message
func Error(v ...interface{}) {
	e := _logger.newEntry(ErrorLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Errorf level formatted message
func Errorf(msg string, v ...interface{}) {
	e := _logger.newEntry(ErrorLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Panic level formatted message
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	e := _logger.newEntry(PanicLevel, s, nil, skipLevel)
	_logger.handleEntry(e)
	panic(s)
}

// Panicf level formatted message
func Panicf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	e := _logger.newEntry(PanicLevel, s, nil, skipLevel)
	_logger.handleEntry(e)
	panic(s)
}

// Fatal level formatted message, followed by an exit.
func Fatal(v ...interface{}) {
	e := _logger.newEntry(FatalLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
	exitFunc(1)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	e := _logger.newEntry(FatalLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	_logger.handleEntry(e)
	exitFunc(1)
}

// WithFields returns a log Entry with fields set
func WithFields(fields Fields) Logger {
	e := _logger.newEntry(InfoLevel, "", fields, skipLevel)
	return e
}

// StackTrace creates a new log Entry with pre-populated field with stack trace.
func StackTrace() Logger {
	trace := make([]byte, 4096)
	runtime.Stack(trace, true)
	customFields := Fields{
		"stack_trace": string(trace) + "\n",
	}
	return _logger.newEntry(DebugLevel, "", customFields, skipLevel)
}

// SetAppID set a constant application key
// that will be set on all log Entry objects
func SetAppID(id string) {
	_logger.appID = id
}

// AppID return an application key
func AppID() string {
	return _logger.appID
}

const (
	requestIdKey = "log_requestIdKey"
)

// NewContext return a new context with a logger value
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, requestIdKey, logger)
}

// FromContext return a logger from the context
func FromContext(ctx context.Context) Logger {
	val, ok := ctx.Value(requestIdKey).(*Entry)
	if ok {
		return val
	}
	return _logger
}
