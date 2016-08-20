package log

import (
	"fmt"
	"sync"
	"time"
)

type Fields map[string]interface{}

type Entry struct {
	wg     *sync.WaitGroup
	logger *Logger

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
func (e *Entry) Warning(v ...interface{}) {
	e.Level = WarningLevel
	e.Message = fmt.Sprint(v...)

	e.logger.handleEntry(e)
}

// Error level message.
func (e *Entry) Error(v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprint(v...)

	e.logger.handleEntry(e)
}

// Fatal level message.
func (e *Entry) Fatal(v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprint(v...)

	e.logger.handleEntry(e)
}

// Consumed lets the Entry and subsequently the Logger
// instance know that it has been used by a handler
func (e *Entry) Consumed() {
	if e.wg != nil {
		e.wg.Done()
	}
}
