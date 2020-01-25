package log

import (
	"context"
	"os"
	"sync"
)

// Logger is the default instance of the log package
var (
	_logger = new()
)

// Handler is an interface that log handlers need to be implemented
type Handler interface {
	Log(Entry)
}

type logger struct {
	host          string
	appID         string
	handlers      map[Level][]Handler
	defaultFields Fields
	rw            sync.RWMutex
}

func new() *logger {
	hostname, _ := os.Hostname()
	logger := logger{
		host:          hostname,
		handlers:      map[Level][]Handler{},
		defaultFields: make(Fields, 0),
	}

	return &logger
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func RegisterHandler(handler Handler, levels ...Level) {
	_logger.rw.Lock()
	defer _logger.rw.Unlock()

	for _, level := range levels {
		_logger.handlers[level] = append(_logger.handlers[level], handler)
	}
}

// Debug level formatted message.
func Debug(msg string) {
	e := newEntry(_logger)
	e.Debug(msg)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Debugf(msg, v...)
}

// Println Info level message.
func Println(msg string) {
	e := newEntry(_logger)
	e.Println(msg)
}

// Print Info level message.
func Print(msg string) {
	e := newEntry(_logger)
	e.Print(msg)
}

// Info level formatted message.
func Info(msg string) {
	e := newEntry(_logger)
	e.Info(msg)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Infof(msg, v...)
}

// Warn level formatted message.
func Warn(msg string) {
	e := newEntry(_logger)
	e.Warn(msg)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Warnf(msg, v...)
}

// Error level formatted message
func Error(msg string) {
	e := newEntry(_logger)
	e.Error(msg)
}

// Errorf level formatted message
func Errorf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Errorf(msg, v...)
}

// Panic level formatted message
func Panic(msg string) {
	e := newEntry(_logger)
	e.Panic(msg)
}

// Panicf level formatted message
func Panicf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Panicf(msg, v...)
}

// Fatal level formatted message, followed by an exit.
func Fatal(msg string) {
	e := newEntry(_logger)
	e.Fatal(msg)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Fatalf(msg, v...)
}

// WithFields returns a log Entry with fields set
func WithFields(fields Fields) Entry {
	e := newEntry(_logger)
	return e.WithFields(fields)
}

// WithField returns a new entry with the `key` and `value` set.
func WithField(key string, value interface{}) Entry {
	e := newEntry(_logger)
	return e.WithField(key, value)
}

// WithDefaultFields adds fields to every entry instance
func WithDefaultFields(fields Fields) {
	newFields := Fields{}

	if _logger.defaultFields != nil {
		for k, val := range _logger.defaultFields {
			newFields[k] = val
		}
	}

	if fields != nil {
		for k, val := range fields {
			newFields[k] = val
		}
	}

	_logger.defaultFields = newFields
}

// WithError returns a new entry with the "error" set to `err`.
func WithError(err error) Entry {
	e := newEntry(_logger)
	return e.WithError((err))
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func Trace(msg string) Entry {
	e := newEntry(_logger)
	return e.Trace(msg)
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

var (
	ctxKey = &struct {
		name string
	}{
		name: "log",
	}
)

// NewContext return a new context with a logger value
func NewContext(ctx context.Context, e Entry) context.Context {
	return context.WithValue(ctx, ctxKey, e)
}

// FromContext return a logger from the context
func FromContext(ctx context.Context) Entry {
	v := ctx.Value(ctxKey)
	if v == nil {
		return newEntry(_logger)
	}

	return v.(Entry)
}
