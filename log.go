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

// Logger interface for logging
type Logger interface {
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
	entryPool        sync.Pool
	channels         LevelHandlerChannels
	callerInfoLevels [6]bool
	appID            string
}

func new() *logger {
	hostname, _ := os.Hostname()
	logger := &logger{
		host:     hostname,
		channels: make(LevelHandlerChannels),
		callerInfoLevels: [6]bool{
			true,
			false,
			false,
			true,
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
	entry.Timestamp = time.Now().UTC()
	return entry
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func (l *logger) RegisterHandler(handler Handler, levels ...Level) {
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
func (l *logger) Debug(v ...interface{}) {
	e := l.newEntry(DebugLevel, fmt.Sprint(v...), nil, skipLevel)
	l.handleEntry(e)
}

// Debugf level formatted message.
func (l *logger) Debugf(msg string, v ...interface{}) {
	e := l.newEntry(DebugLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	l.handleEntry(e)
}

// Info level formatted message.
func (l *logger) Info(v ...interface{}) {
	e := l.newEntry(InfoLevel, fmt.Sprint(v...), nil, skipLevel)
	l.handleEntry(e)
}

// Infof level formatted message.
func (l *logger) Infof(msg string, v ...interface{}) {
	e := l.newEntry(InfoLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	l.handleEntry(e)
}

// Warn level formatted message.
func (l *logger) Warn(v ...interface{}) {
	e := l.newEntry(WarnLevel, fmt.Sprint(v...), nil, skipLevel)
	l.handleEntry(e)
}

// Warnf level formatted message.
func (l *logger) Warnf(msg string, v ...interface{}) {
	e := l.newEntry(WarnLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	l.handleEntry(e)
}

// Error level formatted message.
func (l *logger) Error(v ...interface{}) {
	e := l.newEntry(ErrorLevel, fmt.Sprint(v...), nil, skipLevel)
	l.handleEntry(e)
}

// Errorf level formatted message.
func (l *logger) Errorf(msg string, v ...interface{}) {
	e := l.newEntry(ErrorLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	l.handleEntry(e)
}

// Panic level formatted message.
func (l *logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	e := l.newEntry(PanicLevel, s, nil, skipLevel)
	l.handleEntry(e)
	panic(s)
}

// Panic level formatted message.
func (l *logger) Panicf(msg string, v ...interface{}) {
	s := fmt.Sprintf(msg, v...)
	e := l.newEntry(PanicLevel, s, nil, skipLevel)
	l.handleEntry(e)
	panic(s)
}

// Fatal level formatted message.
func (l *logger) Fatal(v ...interface{}) {
	e := l.newEntry(FatalLevel, fmt.Sprint(v...), nil, skipLevel)
	l.handleEntry(e)
	exitFunc(1)
}

// Fatal level formatted message.
func (l *logger) Fatalf(msg string, v ...interface{}) {
	e := l.newEntry(FatalLevel, fmt.Sprintf(msg, v...), nil, skipLevel)
	l.handleEntry(e)
	exitFunc(1)
}

// WithFields returns a log Entry with fields set
func (l *logger) WithFields(fields Fields) Logger {
	e := l.newEntry(InfoLevel, "", fields, skipLevel)
	return e
}

// SetAppID set a constant application key
// that will be set on all log Entry objects
func (l *logger) SetAppID(id string) {
	l.appID = id
}

// AppID return an application key
func (l *logger) AppID() string {
	return l.appID
}

func (l *logger) handleEntry(e *Entry) {
	if e.Line == 0 && l.callerInfoLevels[e.Level] {
		_, e.File, e.Line, _ = runtime.Caller(e.calldepth)
	}

	channels, ok := l.channels[e.Level]
	if ok {
		e.Lock()
		e.wg.Add(len(channels))
		for _, ch := range channels {
			ch <- e
		}
		e.wg.Wait()
		e.Unlock()
	}

	e.reset()
	l.entryPool.Put(e)
}
