package log

import (
	"fmt"
	"runtime"
	"time"
)

const entryCalldepth = 2

type Fields map[string]interface{}

type Entry struct {
	logger    *logger
	calldepth int

	AppID     string    `json:"app_id"`
	Host      string    `json:"host"`
	Level     Level     `json:"level"`
	Message   string    `json:"message"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Timestamp time.Time `json:"timestamp"`
	Fields    Fields    `json:"fields"`
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func (e *Entry) RegisterHandler(handler Handler, levels ...Level) {
	ch := handler.Run()

	for _, level := range levels {
		channels, ok := e.logger.channels[level]
		if !ok {
			channels = make(HandlerChannels, 0)
		}

		e.logger.channels[level] = append(channels, ch)
	}
}

// Debug level message.
func (e *Entry) Debug(v ...interface{}) {
	e.handler(DebugLevel, fmt.Sprint(v...))
}

// Debugf level message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.handler(DebugLevel, fmt.Sprintf(msg, v...))
}

// Println Info level message.
func (e *Entry) Println(v ...interface{}) {
	e.handler(InfoLevel, fmt.Sprint(v...))
}

// Print Info level message.
func (e *Entry) Print(v ...interface{}) {
	e.handler(InfoLevel, fmt.Sprint(v...))
}

// Info level message.
func (e *Entry) Info(v ...interface{}) {
	e.handler(InfoLevel, fmt.Sprint(v...))
}

// Infof level message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.handler(InfoLevel, fmt.Sprintf(msg, v...))
}

// Warn level message.
func (e *Entry) Warn(v ...interface{}) {
	e.handler(WarnLevel, fmt.Sprint(v...))
}

// Warnf level message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.handler(WarnLevel, fmt.Sprintf(msg, v...))
}

// Error level message.
func (e *Entry) Error(v ...interface{}) {
	e.handler(ErrorLevel, fmt.Sprint(v...))
}

// Errorf level message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.handler(ErrorLevel, fmt.Sprintf(msg, v...))
}

// Panic level message.
func (e *Entry) Panic(v ...interface{}) {
	e.handler(PanicLevel, fmt.Sprint(v...))
	panic(e.Message)
}

// Panicf level message.
func (e *Entry) Panicf(msg string, v ...interface{}) {
	e.handler(PanicLevel, fmt.Sprintf(msg, v...))
	panic(e.Message)
}

// Fatal level message.
func (e *Entry) Fatal(v ...interface{}) {
	e.handler(FatalLevel, fmt.Sprint(v...))
	exitFunc(1)
}

// Fatalf level message.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.handler(FatalLevel, fmt.Sprintf(msg, v...))
	exitFunc(1)
}

// WithFields adds the provided fieldsto the current entry
func (e *Entry) WithFields(fields Fields) Logger {
	data := Fields{}
	if e.Fields != nil {
		for k, val := range e.Fields {
			data[k] = val
		}
	}
	if fields != nil {
		for k, val := range fields {
			data[k] = val
		}
	}
	return e.logger.newEntry(InfoLevel, "", data, entryCalldepth)
}

// Consumed lets the Entry and subsequently the Logger
// instance know that it has been used by a handler
func (e *Entry) Consumed() {
	e.logger.entryPool.Put(e)
}

func (e *Entry) handler(lv Level, msg string) {
	file := ""
	line := 0
	now := time.Now().UTC()
	if e.logger.callerInfoLevels[lv] {
		_, file, line, _ = runtime.Caller(e.calldepth)
	}
	channels, ok := e.logger.channels[lv]
	if ok {
		for _, ch := range channels {
			// new entry here, to prevent race condition
			entry := e.logger.newEntry(lv, msg, e.Fields, entryCalldepth)
			entry.File = file
			entry.Line = line
			entry.Timestamp = now
			ch <- entry
		}
	}
}
