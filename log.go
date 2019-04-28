package log

import (
	"os"
	"strings"
	"sync"
)

// Logger is the default instance of the log package
var (
	exitFunc = os.Exit
)

// Logger interface for logging
type Logger interface {
	Println(v ...interface{})
	Print(v ...interface{})
	Debug(v ...interface{})
	Debugf(msg string, v ...interface{})
	Info(v ...interface{})
	Infof(msg string, v ...interface{})
	Warn(v ...interface{})
	Warnf(msg string, v ...interface{})
	Error(v ...interface{})
	Errorf(msg string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(msg string, v ...interface{})
	Panic(v ...interface{})
	Panicf(msg string, v ...interface{})
	WithFields(fields Fields) Logger
}

// HandlerChannels is an array of handler channels
type HandlerChannels []chan<- *Entry

// LevelHandlerChannels is a group of Handler channels mapped by Level
type LevelHandlerChannels map[Level]HandlerChannels

type Handler interface {
	Run() chan<- *Entry
}

type logger struct {
	host             string
	entryPool        *sync.Pool
	channels         LevelHandlerChannels
	callerInfoLevels [6]bool
	appID            string
}

func new() *Entry {
	hostname, _ := os.Hostname()
	logger := &logger{
		host:     hostname,
		channels: make(LevelHandlerChannels),
		callerInfoLevels: [6]bool{
			true,  // debug
			false, // info
			false, // warn
			true,  // error
			true,  // panic
			true,  // fatal
		},
	}

	logger.entryPool = &sync.Pool{
		New: func() interface{} {
			return &Entry{
				logger: logger,
			}
		},
	}

	return logger.entryPool.Get().(*Entry)
}

func (l *logger) newEntry(level Level, message string, fields Fields, calldepth int) *Entry {
	entry := l.entryPool.Get().(*Entry)
	entry.logger = l
	entry.calldepth = calldepth
	entry.Host = l.host
	entry.AppID = l.appID
	entry.Line = 0
	entry.File = ""
	entry.Level = level
	entry.Message = strings.TrimSpace(message)
	entry.Fields = fields
	return entry
}
