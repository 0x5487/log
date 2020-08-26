package log

import (
	"context"
	"fmt"
	"time"
)

// Context use for meta data
type Context struct {
	logger *logger
	buf    []byte
}

func newContext(l *logger) Context {
	c := Context{
		logger: l,
	}

	if len(l.buf) > 0 {
		c.buf = copyBytes(l.buf)
	} else {
		c.buf = make([]byte, 0, 500)
		c.buf = enc.AppendBeginMarker(c.buf)
	}

	return c
}

func copyBytes(src []byte) []byte {
	newBuf := make([]byte, len(src))
	if len(src) > 0 {
		copy(newBuf, src)
	}
	return newBuf
}

// SaveToDefault save the current context to default logger and these context to be printed with every entry
func (c Context) SaveToDefault() {
	c.logger.rwMutex.Lock()
	defer c.logger.rwMutex.Unlock()

	c.logger.buf = copyBytes(c.buf)
}

// Debug level formatted message.
func (c Context) Debug(msg string) {
	e := newEntry(_logger, c.buf)
	e.Debug(msg)
}

// Debugf level formatted message.
func (c Context) Debugf(msg string, v ...interface{}) {
	e := newEntry(_logger, c.buf)
	e.Debugf(msg, v...)
}

// Info level formatted message.
func (c Context) Info(msg string) {
	e := newEntry(_logger, c.buf)
	e.Info(msg)
}

// Infof level formatted message.
func (c Context) Infof(msg string, v ...interface{}) {
	e := newEntry(_logger, c.buf)
	e.Infof(msg, v...)
}

// Warn level formatted message.
func (c Context) Warn(msg string) {
	e := newEntry(_logger, c.buf)
	e.Warn(msg)
}

// Warnf level formatted message.
func (c Context) Warnf(msg string, v ...interface{}) {
	e := newEntry(_logger, c.buf)
	e.Warnf(msg, v...)
}

// Error level formatted message
func (c Context) Error(msg string) {
	e := newEntry(_logger, c.buf)
	e.Error(msg)
}

// Errorf level formatted message
func (c Context) Errorf(msg string, v ...interface{}) {
	e := newEntry(_logger, c.buf)
	e.Errorf(msg, v...)
}

// Panic level formatted message
func (c Context) Panic(msg string) {
	e := newEntry(_logger, c.buf)
	e.Panic(msg)
}

// Panicf level formatted message
func (c Context) Panicf(msg string, v ...interface{}) {
	e := newEntry(_logger, c.buf)
	e.Panicf(msg, v...)
}

// Fatal level formatted message
func (c Context) Fatal(msg string) {
	e := newEntry(_logger, c.buf)
	e.Fatal(msg)
}

// Fatalf level formatted message
func (c Context) Fatalf(msg string, v ...interface{}) {
	e := newEntry(_logger, c.buf)
	e.Fatalf(msg, v...)
}

// Str add string field to current context
func (c Context) Str(key string, val string) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendString(c.buf, val)
	return c
}

// Strs add string field to current context
func (c Context) Strs(key string, val []string) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendStrings(c.buf, val)
	return c
}

// Bool add bool field to current context
func (c Context) Bool(key string, val bool) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendBool(c.buf, val)
	return c
}

// Int add Int field to current context
func (c Context) Int(key string, val int) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendInt(c.buf, val)
	return c
}

// Ints add Int field to current context
func (c Context) Ints(key string, val []int) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendInts(c.buf, val)
	return c
}

// Int8 add Int8 field to current context
func (c Context) Int8(key string, val int8) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendInt8(c.buf, val)
	return c
}

// Int16 add Int16 field to current context
func (c Context) Int16(key string, val int16) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendInt16(c.buf, val)
	return c
}

// Int32 add Int32 field to current context
func (c Context) Int32(key string, val int32) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendInt32(c.buf, val)
	return c
}

// Int64 add Int64 field to current context
func (c Context) Int64(key string, val int64) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendInt64(c.buf, val)
	return c
}

// Uint add Uint field to current context
func (c Context) Uint(key string, val uint) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendUint(c.buf, val)
	return c
}

// Uint8 add Uint8 field to current context
func (c Context) Uint8(key string, val uint8) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendUint8(c.buf, val)
	return c
}

// Uint16 add Uint16 field to current context
func (c Context) Uint16(key string, val uint16) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendUint16(c.buf, val)
	return c
}

// Uint32 add Uint32 field to current context
func (c Context) Uint32(key string, val uint32) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendUint32(c.buf, val)
	return c
}

// Uint64 add Uint64 field to current context
func (c Context) Uint64(key string, val uint64) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendUint64(c.buf, val)
	return c
}

// Float32 add float32 field to current context
func (c Context) Float32(key string, val float32) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendFloat32(c.buf, val)
	return c
}

// Float64 add Float64 field to current context
func (c Context) Float64(key string, val float64) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendFloat64(c.buf, val)
	return c
}

// Err add error field to current context
func (c Context) Err(err error) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, "error")
	c.buf = enc.AppendString(c.buf, fmt.Sprintf("%+v", err))
	return c
}

// StackTrace adds stack_trace field to the current context
func (c Context) StackTrace() Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, "stack_trace")
	c.buf = enc.AppendString(c.buf, getStackTrace())
	return c
}

// Time adds Time field to current context
func (c Context) Time(key string, val time.Time) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendTime(c.buf, val, time.RFC3339)
	return c
}

// Times adds Time field to current context
func (c Context) Times(key string, val []time.Time) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendTimes(c.buf, val, time.RFC3339)
	return c
}

// Interface adds the field key with i marshaled using reflection.
func (c Context) Interface(key string, val interface{}) Context {
	c.buf = copyBytes(c.buf)
	c.buf = enc.AppendKey(c.buf, key)
	c.buf = enc.AppendInterface(c.buf, val)
	return c
}

// WithContext return a new context with a log context value
func (c Context) WithContext(ctx context.Context) context.Context {
	return newStdContext(ctx, c)
}
