package log

import (
	"context"
	"fmt"
	"runtime"
)

var entry = new()

// SetLogger replace default log
// Example:
//   log.SetLogger(log.WithFields(log.Fields{"location": os.Getenv("LOCATION")}))
func SetLogger(e Logger) {
	entry = e.(*Entry)
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func RegisterHandler(handler Handler, levels ...Level) {
	entry.RegisterHandler(handler, levels...)
}

// Debug level formatted message.
func Debug(v ...interface{}) {
	entry.handler(DebugLevel, fmt.Sprint(v...))
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	entry.handler(DebugLevel, fmt.Sprintf(msg, v...))
}

// Println Info level message.
func Println(v ...interface{}) {
	entry.handler(InfoLevel, fmt.Sprint(v...))
}

// Print Info level message.
func Print(v ...interface{}) {
	entry.handler(InfoLevel, fmt.Sprint(v...))
}

// Info level formatted message.
func Info(v ...interface{}) {
	entry.handler(InfoLevel, fmt.Sprint(v...))
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	entry.handler(InfoLevel, fmt.Sprintf(msg, v...))
}

// Warn level formatted message.
func Warn(v ...interface{}) {
	entry.handler(WarnLevel, fmt.Sprint(v...))
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	entry.handler(WarnLevel, fmt.Sprintf(msg, v...))
}

// Error level formatted message
func Error(v ...interface{}) {
	entry.handler(ErrorLevel, fmt.Sprint(v...))
}

// Errorf level formatted message
func Errorf(msg string, v ...interface{}) {
	entry.handler(ErrorLevel, fmt.Sprintf(msg, v...))
}

// Panic level formatted message
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	entry.handler(PanicLevel, s)
	panic(s)
}

// Panicf level formatted message
func Panicf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	entry.handler(PanicLevel, s)
	panic(s)
}

// Fatal level formatted message, followed by an exit.
func Fatal(v ...interface{}) {
	entry.handler(FatalLevel, fmt.Sprint(v...))
	exitFunc(1)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	entry.handler(FatalLevel, fmt.Sprintf(msg, v...))
	exitFunc(1)
}

// WithFields returns a log Entry with fields set
func WithFields(fields Fields) Logger {
	return entry.WithFields(fields)
}

// StackTrace creates a new log Entry with pre-populated field with stack trace.
func StackTrace() *Entry {
	trace := make([]byte, 4096)
	runtime.Stack(trace, true)
	customFields := Fields{
		"stack_trace": string(trace) + "\n",
	}
	return entry.logger.newEntry(ErrorLevel, "", customFields, entryCalldepth)
}

// SetAppID set a constant application key
// that will be set on all log Entry objects
func SetAppID(id string) {
	entry.logger.appID = id
}

// AppID return an application key
func AppID() string {
	return entry.logger.appID
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
	return entry
}
