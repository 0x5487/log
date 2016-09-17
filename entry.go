package log

import (
	"fmt"
	"sync"
	"time"
)

type Fields map[string]interface{}

type Entry struct {
	wg        *sync.WaitGroup
	logger    *logger
	calldepth int

	AppID     string
	Host      string
	Level     Level  `json:"level"`
	Message   string `json:"message"`
	File      string
	Line      int       `json:"line"`
	Timestamp time.Time `json:"timestamp"`
	Fields    Fields
}

// Debug level message.
func (e *Entry) Debug(v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprint(v...)

	e.logger.handleEntry(e)
}

// Info level message.
func (e *Entry) Info(v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprint(v...)

	e.logger.handleEntry(e)
}

// Warning level message.
func (e *Entry) Warn(v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
}

// Error level message.
func (e *Entry) Error(v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprint(v...)
	e.logger.handleEntry(e)
}

// Panic level message.
func (e *Entry) Panic(v ...interface{}) {
	e.Level = PanicLevel
	e.Message = fmt.Sprint(v...)
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

// Consumed lets the Entry and subsequently the Logger
// instance know that it has been used by a handler
func (e *Entry) Consumed() {
	if e.wg != nil {
		e.wg.Done()
	}
}
