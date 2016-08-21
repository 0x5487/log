package log

// Level of the log
type Level uint8

// Log levels.
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// AllLevels is an array of all log levels, for easier registering of all levels to a handler
var AllLevels = []Level{
	DebugLevel,
	InfoLevel,
	WarnLevel,
	ErrorLevel,
	FatalLevel,
}

var levelNames = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// String returns the string representation of a logging level.
func (p Level) String() string {
	return levelNames[p]
}
