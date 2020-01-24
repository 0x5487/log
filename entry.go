package log

import (
	"fmt"
	"time"
)

type Fields map[string]interface{}

// Entry defines a single log entry
type Entry struct {
	logger *logger

	AppID     string    `json:"app_id"`
	Host      string    `json:"host"`
	Level     Level     `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Fields    Fields    `json:"fields"`
}

func newEntry(l *logger) Entry {
	e := Entry{}
	e.logger = l
	e.Host = l.host
	e.AppID = l.appID
	e.Fields = l.defaultFields
	return e
}

// Debug level message.
func (e Entry) Debug(v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
}

// Debugf level message.
func (e Entry) Debugf(msg string, v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Println Info level message.
func (e Entry) Println(v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
}

// Print Info level message.
func (e Entry) Print(v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
}

// Info level message.
func (e Entry) Info(v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
}

// Infof level message.
func (e Entry) Infof(msg string, v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Warn level message.
func (e Entry) Warn(v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
}

// Warnf level message.
func (e Entry) Warnf(msg string, v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Error level message.
func (e Entry) Error(v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
}

// Errorf level message.
func (e Entry) Errorf(msg string, v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Panic level message.
func (e Entry) Panic(v ...interface{}) {
	e.Level = PanicLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
	exitFunc(1)
}

// Panicf level message.
func (e Entry) Panicf(msg string, v ...interface{}) {
	e.Level = PanicLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
	exitFunc(1)
}

// Fatal level message.
func (e Entry) Fatal(v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprint(v...)
	handler(e)
	exitFunc(1)
}

// Fatalf level message.
func (e Entry) Fatalf(msg string, v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
	exitFunc(1)
}

// WithFields adds the provided fields to the current entry
func (e Entry) WithFields(fields Fields) Entry {
	newFields := Fields{}
	if e.Fields != nil {
		for k, val := range e.Fields {
			newFields[k] = val
		}
	}
	if fields != nil {
		for k, val := range fields {
			newFields[k] = val
		}
	}

	e.Fields = newFields
	return e
}

func handler(e Entry) {
	e.Timestamp = time.Now().UTC()

	e.logger.rw.RLock()
	defer e.logger.rw.RUnlock()

	for _, h := range e.logger.handlers[e.Level] {
		h.Log(e)
	}
}
