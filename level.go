package log

import "strings"

// Level of the log
type Level uint8

// Log levels.
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

// AllLevels is an array of all log levels, for easier registering of all levels to a handler
var AllLevels = []Level{
	DebugLevel,
	InfoLevel,
	WarnLevel,
	ErrorLevel,
	PanicLevel,
	FatalLevel,
}

var levelNames = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"PANIC",
	"FATAL",
}

// String returns the string representation of a logging level.
func (p Level) String() string {
	return levelNames[p]
}

func GetLevelsFromMinLevel(minLevel string) []Level {
	minLevel = strings.ToLower(minLevel)
	switch minLevel {
	case "debug":
		return AllLevels
	case "info":
		return []Level{
			InfoLevel,
			WarnLevel,
			ErrorLevel,
			PanicLevel,
			FatalLevel,
		}
	case "warn":
		return []Level{
			WarnLevel,
			ErrorLevel,
			PanicLevel,
			FatalLevel,
		}
	case "error":
		return []Level{
			ErrorLevel,
			PanicLevel,
			FatalLevel,
		}
	case "panic":
		return []Level{
			PanicLevel,
			FatalLevel,
		}
	case "fatal":
		return []Level{
			FatalLevel,
		}
	default:
		return AllLevels
	}
}
