package log

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Level of the log
type Level uint8

// Logger is the default instance of the log package
var (
	exitFunc  = os.Exit
	skipLevel = 2
)

// Log levels.
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	AlertLevel
	FatalLevel
)

// HandlerChannels is an array of handler channels
type HandlerChannels []chan<- *Entry

// LevelHandlerChannels is a group of Handler channels mapped by Level
type LevelHandlerChannels map[Level]HandlerChannels

type Handler interface {
	Run() chan<- *Entry
}

type Logger struct {
	entryPool *sync.Pool
	channels  LevelHandlerChannels
	Name      string
	AppID     string
}

func New(appID string) *Logger {
	return &Logger{
		channels: make(LevelHandlerChannels),
		AppID:    appID,
	}
}

func (l *Logger) newEntry(level Level, message string, fields CustomFields, calldepth int) *Entry {
	entry := l.entryPool.Get().(*Entry)
	entry.Line = 0
	entry.File = entry.File[0:0]
	entry.calldepth = calldepth
	entry.Level = level
	entry.Message = strings.TrimRight(message, cutset) // need to trim for adding fields later in handlers + why send uneeded whitespace
	entry.CustomFields = fields
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

func (l *Logger) Info(v ...interface{}) {
	e := l.newEntry(InfoLevel, fmt.Sprint(v...), nil, skipLevel)
	l.HandleEntry(e)
}

func (l *Logger) HandleEntry(e *Entry) {
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
