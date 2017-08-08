package log

import (
	"fmt"
	"sync"
	"time"
)

type Fields map[string]interface{}

type Entry struct {
	sync.RWMutex
	wg        *sync.WaitGroup
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

// Debug level message.
func (e *Entry) Debug(v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
}

// Debugf level message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprintf(msg, v...)
	e.logger.handleEntry(e)
}

// Info level message.
func (e *Entry) Info(v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
}

// Infof level message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprintf(msg, v...)
	e.logger.handleEntry(e)
}

// Warn level message.
func (e *Entry) Warn(v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
}

// Warnf level message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprintf(msg, v...)
	e.logger.handleEntry(e)
}

// Error level message.
func (e *Entry) Error(v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
}

// Errorf level message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprintf(msg, v...)
	e.logger.handleEntry(e)
}

// Panic level message.
func (e *Entry) Panic(v ...interface{}) {
	e.Level = PanicLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
	panic(e.Message)
}

// Panicf level message.
func (e *Entry) Panicf(msg string, v ...interface{}) {
	e.Level = PanicLevel
	e.Message = fmt.Sprintf(msg, v...)
	e.logger.handleEntry(e)
	panic(e.Message)
}

// Fatal level message.
func (e *Entry) Fatal(v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
	exitFunc(1)
}

// Fatalf level message.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprintf(msg, v...)
	e.logger.handleEntry(e)
	exitFunc(1)
}

// WithFields adds the provided fieldsto the current entry
func (e *Entry) WithFields(fields Fields) Logger {
	if e.Fields == nil {
		e.Fields = fields
	} else {
		for k, val := range fields {
			e.Fields[k] = val
		}
	}
	return e
}

// Consumed lets the Entry and subsequently the Logger
// instance know that it has been used by a handler
func (e *Entry) Consumed() {
	if e.wg != nil {
		e.wg.Done()
	}
}

func (e *Entry) reset() {
	e.File = ""
	e.Line = 0
}
