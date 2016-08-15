package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger is the default instance of the log package
var (
	exitFunc  = os.Exit
	skipLevel = 2
)

// HandlerChannels is an array of handler channels
type HandlerChannels []chan<- *Entry

// LevelHandlerChannels is a group of Handler channels mapped by Level
type LevelHandlerChannels map[Level]HandlerChannels

type Handler interface {
	Run() chan<- *Entry
}

type Logger struct {
	entryPool        sync.Pool
	channels         LevelHandlerChannels
	callerInfoLevels [5]bool
	AppID            string
}

func New() *Logger {
	logger := &Logger{
		channels: make(LevelHandlerChannels),
		callerInfoLevels: [5]bool{
			true,
			false,
			false,
			true,
			true,
		},
	}

	logger.entryPool.New = func() interface{} {
		return &Entry{
			wg: &sync.WaitGroup{},
		}
	}

	return logger
}

func (l *Logger) newEntry(level Level, message string, fields Fields) *Entry {
	entry := l.entryPool.Get().(*Entry)
	entry.logger = l
	entry.Line = 0
	entry.File = ""
	entry.Level = level
	entry.Message = strings.TrimSpace(message)
	entry.Fields = fields
	entry.Timestamp = time.Now().UTC()

	return entry
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func (l *Logger) RegisterHandler(handler Handler, levels ...Level) {
	ch := handler.Run()

	for _, level := range levels {
		channels, ok := l.channels[level]
		if !ok {
			channels = make(HandlerChannels, 0)
		}

		l.channels[level] = append(channels, ch)
	}
}

// Debug level formatted message.
func (l *Logger) Debug(v ...interface{}) {
	e := l.newEntry(DebugLevel, fmt.Sprint(v...), nil)
	l.handleEntry(e)
}

// Info level formatted message.
func (l *Logger) Info(v ...interface{}) {
	e := l.newEntry(InfoLevel, fmt.Sprint(v...), nil)
	l.handleEntry(e)
}

// Warning level formatted message.
func (l *Logger) Warning(v ...interface{}) {
	e := l.newEntry(WarningLevel, fmt.Sprint(v...), nil)
	l.handleEntry(e)
}

// Error level formatted message, followed by an exit.
func (l *Logger) Error(v ...interface{}) {
	e := l.newEntry(ErrorLevel, fmt.Sprint(v...), nil)
	l.handleEntry(e)
}

// Fatal level formatted message, followed by an exit.
func (l *Logger) Fatal(v ...interface{}) {
	e := l.newEntry(FatalLevel, fmt.Sprint(v...), nil)
	l.handleEntry(e)
	exitFunc(1)
}

// WithFields returns a log Entry with fields set
func (l *Logger) WithFields(fields Fields) *Entry {
	e := l.newEntry(InfoLevel, "", fields)
	l.handleEntry(e)
	return e
}

func (l *Logger) handleEntry(e *Entry) {
	if e.Line == 0 && l.callerInfoLevels[e.Level] {
		_, e.File, e.Line, _ = runtime.Caller(2)
	}

	channels, ok := l.channels[e.Level]

	if ok {
		e.wg.Add(len(channels))
		for _, ch := range channels {
			ch <- e
		}
		e.wg.Wait()
	}

	l.entryPool.Put(e)
}
