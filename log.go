package log

import (
	"context"
	"fmt"
	stdlog "log"
	"runtime"
	"strings"
	"sync"
)

// Logger is the default instance of the log package
var (
	_logger = new()

	// ErrorHandler is called whenever handler fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)

	// AutoStaceTrace add stack trace into entry when use `Error`, `Panic`, `Fatal` level.
	// Default: true
	AutoStaceTrace = true
)

// Handler is an interface that log handlers need to be implemented
type Handler interface {
	BeforeWriting(*Entry) error
	Write([]byte) error
}

// Flusher is an interface that allow handles have the ability to clear buffer and close connection
type Flusher interface {
	Flush() error
}

// Hookfunc is an func that allow us to do something before writing
type Hookfunc func(*Entry) error

type logger struct {
	handles              []Handler
	hooks                []Hookfunc
	leveledHandlers      map[Level][]Handler
	cacheLeveledHandlers func(level Level) []Handler
	rwMutex              sync.RWMutex
	buf                  []byte
}

func new() *logger {
	logger := logger{
		leveledHandlers: map[Level][]Handler{},
	}

	logger.cacheLeveledHandlers = logger.getLeveledHandlers()
	return &logger
}

func (l *logger) getLeveledHandlers() func(level Level) []Handler {
	debugHandlers := l.leveledHandlers[DebugLevel]
	infoHandlers := l.leveledHandlers[InfoLevel]
	warnHandlers := l.leveledHandlers[WarnLevel]
	errorHandlers := l.leveledHandlers[ErrorLevel]
	panicHandlers := l.leveledHandlers[PanicLevel]
	fatalHandlers := l.leveledHandlers[FatalLevel]

	return func(level Level) []Handler {
		switch level {
		case DebugLevel:
			return debugHandlers
		case InfoLevel:
			return infoHandlers
		case WarnLevel:
			return warnHandlers
		case ErrorLevel:
			return errorHandlers
		case PanicLevel:
			return panicHandlers
		case FatalLevel:
			return fatalHandlers
		}

		return []Handler{}
	}
}

// AddHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func AddHandler(handler Handler, levels ...Level) {
	_logger.rwMutex.Lock()
	defer _logger.rwMutex.Unlock()

	for _, level := range levels {
		_logger.leveledHandlers[level] = append(_logger.leveledHandlers[level], handler)
	}

	_logger.handles = append(_logger.handles, handler)
	_logger.cacheLeveledHandlers = _logger.getLeveledHandlers()
}

// RemoveAllHandlers removes all handlers
func RemoveAllHandlers() {
	_logger.rwMutex.Lock()
	defer _logger.rwMutex.Unlock()

	_logger.leveledHandlers = map[Level][]Handler{}
	_logger.handles = []Handler{}
	_logger.hooks = []Hookfunc{}
	_logger.cacheLeveledHandlers = _logger.getLeveledHandlers()
}

// AddHook adds a new Hook to log entry
func AddHook(hook Hookfunc) error {
	_logger.rwMutex.Lock()
	defer _logger.rwMutex.Unlock()

	_logger.hooks = append(_logger.hooks, hook)
	return nil
}

// Debug level formatted message
func Debug(msg string) {
	e := newEntry(_logger, _logger.buf)
	e.Debug(msg)
}

// Debugf level formatted message
func Debugf(msg string, v ...interface{}) {
	e := newEntry(_logger, _logger.buf)
	e.Debugf(msg, v...)
}

// Info level formatted message
func Info(msg string) {
	e := newEntry(_logger, _logger.buf)
	e.Info(msg)
}

// Infof level formatted message
func Infof(msg string, v ...interface{}) {
	e := newEntry(_logger, _logger.buf)
	e.Infof(msg, v...)
}

// Warn level formatted message
func Warn(msg string) {
	e := newEntry(_logger, _logger.buf)
	e.Warn(msg)
}

// Warnf level formatted message
func Warnf(msg string, v ...interface{}) {
	e := newEntry(_logger, _logger.buf)
	e.Warnf(msg, v...)
}

// Error level formatted message
func Error(msg string) {
	e := newEntry(_logger, _logger.buf)
	e.Error(msg)
}

// Errorf level formatted message
func Errorf(msg string, v ...interface{}) {
	e := newEntry(_logger, _logger.buf)
	e.Errorf(msg, v...)
}

// Panic level formatted message
func Panic(msg string) {
	e := newEntry(_logger, _logger.buf)
	e.Panic(msg)
}

// Panicf level formatted message
func Panicf(msg string, v ...interface{}) {
	e := newEntry(_logger, nil)
	e.Panicf(msg, v...)
}

// Fatal level formatted message, followed by an exit.
func Fatal(msg string) {
	e := newEntry(_logger, _logger.buf)
	e.Fatal(msg)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	e := newEntry(_logger, _logger.buf)
	e.Fatalf(msg, v...)
}

// Str add string field to current context
func Str(key string, val string) Context {
	c := newContext(_logger)
	return c.Str(key, val)
}

// Bool add bool field to current context
func Bool(key string, val bool) Context {
	c := newContext(_logger)
	return c.Bool(key, val)
}

// Int add Int field to current context
func Int(key string, val int) Context {
	c := newContext(_logger)
	return c.Int(key, val)
}

// Int8 add Int8 field to current context
func Int8(key string, val int8) Context {
	c := newContext(_logger)
	return c.Int8(key, val)
}

// Int16 add Int16 field to current context
func Int16(key string, val int16) Context {
	c := newContext(_logger)
	return c.Int16(key, val)
}

// Int32 add Int32 field to current context
func Int32(key string, val int32) Context {
	c := newContext(_logger)
	return c.Int32(key, val)
}

// Int64 add Int64 field to current context
func Int64(key string, val int64) Context {
	c := newContext(_logger)
	return c.Int64(key, val)
}

// Uint add Uint field to current context
func Uint(key string, val uint) Context {
	c := newContext(_logger)
	return c.Uint(key, val)
}

// Uint8 add Uint8 field to current context
func Uint8(key string, val uint8) Context {
	c := newContext(_logger)
	return c.Uint8(key, val)
}

// Uint16 add Uint16 field to current context
func Uint16(key string, val uint16) Context {
	c := newContext(_logger)
	return c.Uint16(key, val)
}

// Uint32 add Uint32 field to current context
func Uint32(key string, val uint32) Context {
	c := newContext(_logger)
	return c.Uint32(key, val)
}

// Uint64 add Uint64 field to current context
func Uint64(key string, val uint64) Context {
	c := newContext(_logger)
	return c.Uint64(key, val)
}

// Float32 add float32 field to current context
func Float32(key string, val float32) Context {
	c := newContext(_logger)
	return c.Float32(key, val)
}

// Float64 add Float64 field to current context
func Float64(key string, val float64) Context {
	c := newContext(_logger)
	return c.Float64(key, val)
}

// Err add error field to current context
func Err(err error) Context {
	c := newContext(_logger)
	return c.Err(err)
}

// Flush clear all handler's buffer
func Flush() {
	for _, h := range _logger.handles {
		flusher, ok := h.(Flusher)
		if ok {
			err := flusher.Flush()
			if err != nil {
				stdlog.Printf("log: flush log handler: %v", err)
			}
		}
	}
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func Trace(msg string) *Entry {
	e := newEntry(_logger, nil)
	return e.Trace(msg)
}

var (
	ctxKey = &struct {
		name string
	}{
		name: "log",
	}
)

// newStdContext return a new context with a log context value
func newStdContext(ctx context.Context, c Context) context.Context {
	return context.WithValue(ctx, ctxKey, c)
}

// FromContext return a log context from the standard context
func FromContext(ctx context.Context) Context {
	v := ctx.Value(ctxKey)
	if v == nil {
		return newContext(_logger)
	}

	return v.(Context)
}

func getStackTrace() string {
	stackBuf := make([]uintptr, 50)
	length := runtime.Callers(3, stackBuf[:])
	stack := stackBuf[:length]

	var b strings.Builder
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			_, _ = b.WriteString(fmt.Sprintf("\n\tFile: %s, Line: %d. Function: %s", frame.File, frame.Line, frame.Function))
		}
		if !more {
			break
		}
	}
	return b.String()
}
