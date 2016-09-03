package log

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
	_logger.Debug(v...)
}

// Info level formatted message.
func Info(v ...interface{}) {
	_logger.Info(v...)
}

// Warn level formatted message.
func Warn(v ...interface{}) {
	_logger.Warn(v...)
}

// Error level formatted message
func Error(v ...interface{}) {
	_logger.Error(v...)
}

// Panic level formatted message
func Panic(v ...interface{}) {
	_logger.Panic(v...)
}

// Fatal level formatted message, followed by an exit.
func Fatal(v ...interface{}) {
	_logger.Fatal(v...)
}

// WithFields returns a log Entry with fields set
func WithFields(fields Fields) *Entry {
	return _logger.WithFields(fields)
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
