package log

import "fmt"

var (
	_logger *Logger
)

func init() {
	_logger = New()
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

// Info level formatted message.
func Info(v ...interface{}) {
	e := _logger.newEntry(InfoLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Warn level formatted message.
func Warn(v ...interface{}) {
	e := _logger.newEntry(WarnLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Error level formatted message
func Error(v ...interface{}) {
	e := _logger.newEntry(ErrorLevel, fmt.Sprint(v...), nil, skipLevel)
	_logger.handleEntry(e)
}

// Panic level formatted message
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
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

// WithFields returns a log Entry with fields set
func WithFields(fields Fields) *Entry {
	e := _logger.newEntry(InfoLevel, "", fields, skipLevel)
	return e
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