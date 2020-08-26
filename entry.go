package log

import (
	"fmt"
	stdlog "log"
	"os"
	"sync"
	"time"

	"github.com/jasonsoft/log/v2/internal/json"
)

var entryPool = &sync.Pool{
	New: func() interface{} {
		return &Entry{
			buf: make([]byte, 0, 500),
		}
	},
}

var enc = json.Encoder{}

// Entry defines a single log entry
type Entry struct {
	logger *logger
	start  time.Time
	buf    []byte

	Level   Level  `json:"level"`
	Message string `json:"message"`
}

func newEntry(l *logger, buf []byte) *Entry {
	e := entryPool.Get().(*Entry)
	e.logger = l

	if buf == nil {
		e.buf = e.buf[:0]
		e.buf = enc.AppendBeginMarker(e.buf)
	} else {
		e.buf = buf // race condition here because it use context's buf.  However, we create new buf within handler func
	}

	return e
}

func putEntry(e *Entry) {
	// Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16 // 64KiB
	if cap(e.buf) > maxSize {
		return
	}
	entryPool.Put(e)
}

func copyEntry(e *Entry) *Entry {
	newEntry := entryPool.Get().(*Entry)

	if len(e.buf) > cap(newEntry.buf) {
		// append will auto increase slice's capacity  when needed
		newEntry.buf = e.buf[:0]
		newEntry.buf = append(newEntry.buf, e.buf...)
	} else {
		// Copy returns the number of elements copied, which will be the minimum of len(src) and len(dst).
		// https://stackoverflow.com/questions/30182538/why-cant-i-duplicate-a-slice-with-copy
		newEntry.buf = newEntry.buf[:len(e.buf)]
		copy(newEntry.buf, e.buf)
	}

	newEntry.logger = e.logger
	newEntry.start = e.start
	newEntry.Level = e.Level
	newEntry.Message = e.Message
	return newEntry
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (e *Entry) Trace(msg string) *Entry {
	e.Level = InfoLevel
	e.Message = msg
	e.start = time.Now().UTC()
	return e
}

// Stop should be used with Trace, to fire off the completion message. When
// an `err` is passed the "error" field is set, and the log level is error.
func (e *Entry) Stop() {
	e = e.Dur("duration", time.Since(e.start))
	handler(e)
}

// Debug level message.
func (e *Entry) Debug(msg string) {
	e.Level = DebugLevel
	e.Message = msg
	handler(e)
}

// Debugf level message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Info level message.
func (e *Entry) Info(msg string) {
	e.Level = InfoLevel
	e.Message = msg

	handler(e)

}

// Infof level message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Warn level message.
func (e *Entry) Warn(msg string) {
	e.Level = WarnLevel
	e.Message = msg
	handler(e)
}

// Warnf level message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Error level message.
func (e *Entry) Error(msg string) {
	e.Level = ErrorLevel
	e.Message = msg

	if AutoStaceTrace {
		e = e.StackTrace()
	}

	handler(e)
}

// Errorf level message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprintf(msg, v...)

	if AutoStaceTrace {
		e = e.StackTrace()
	}

	handler(e)
}

// Panic level message.
func (e *Entry) Panic(msg string) {
	e.Level = PanicLevel
	e.Message = msg

	if AutoStaceTrace {
		e = e.StackTrace()
	}

	handler(e)
	panic(msg)
}

// Panicf level message.
func (e *Entry) Panicf(msg string, v ...interface{}) {
	e.Level = PanicLevel
	e.Message = fmt.Sprintf(msg, v...)

	if AutoStaceTrace {
		e = e.StackTrace()
	}

	handler(e)
	panic(msg)
}

// Fatal level message.
func (e *Entry) Fatal(msg string) {
	e.Level = FatalLevel
	e.Message = msg

	if AutoStaceTrace {
		e = e.StackTrace()
	}
	handler(e)
	os.Exit(1)
}

// Fatalf level message.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprintf(msg, v...)

	if AutoStaceTrace {
		e = e.StackTrace()
	}
	handler(e)
	os.Exit(1)
}

// Str add string field to current entry
func (e *Entry) Str(key string, val string) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendString(e.buf, val)
	return e
}

// Strs add string field to current entry
func (e *Entry) Strs(key string, val []string) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendStrings(e.buf, val)
	return e
}

// Bool add bool field to current entry
func (e *Entry) Bool(key string, val bool) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendBool(e.buf, val)
	return e
}

// Int adds Int field to current entry
func (e *Entry) Int(key string, val int) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendInt(e.buf, val)
	return e
}

// Ints adds Int field to current entry
func (e *Entry) Ints(key string, val []int) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendInts(e.buf, val)
	return e
}

// Int8 add Int8 field to current entry
func (e *Entry) Int8(key string, val int8) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendInt8(e.buf, val)
	return e
}

// Int16 add Int16 field to current entry
func (e *Entry) Int16(key string, val int16) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendInt16(e.buf, val)
	return e
}

// Int32 adds Int32 field to current entry
func (e *Entry) Int32(key string, val int32) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendInt32(e.buf, val)
	return e
}

// Int64 add Int64 field to current entry
func (e *Entry) Int64(key string, val int64) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendInt64(e.buf, val)
	return e
}

// Uint add Uint field to current entry
func (e *Entry) Uint(key string, val uint) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendUint(e.buf, val)
	return e
}

// Uint8 add Uint8 field to current entry
func (e *Entry) Uint8(key string, val uint8) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendUint8(e.buf, val)
	return e
}

// Uint16 add Uint16 field to current entry
func (e *Entry) Uint16(key string, val uint16) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendUint16(e.buf, val)
	return e
}

// Uint32 add Uint32 field to current entry
func (e *Entry) Uint32(key string, val uint32) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendUint32(e.buf, val)
	return e
}

// Uint64 add Uint64 field to current entry
func (e *Entry) Uint64(key string, val uint64) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendUint64(e.buf, val)
	return e
}

// Float32 add Float32 field to current entry
func (e *Entry) Float32(key string, val float32) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendFloat32(e.buf, val)
	return e
}

// Float64 adds Float64 field to current entry
func (e *Entry) Float64(key string, val float64) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendFloat64(e.buf, val)
	return e
}

// Time adds Time field to current entry
func (e *Entry) Time(key string, val time.Time) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendTime(e.buf, val, time.RFC3339)
	return e
}

// Times adds Time field to current entry
func (e *Entry) Times(key string, val []time.Time) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendTimes(e.buf, val, time.RFC3339)
	return e
}

// Dur adds Duration field to current entry
func (e *Entry) Dur(key string, d time.Duration) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendDuration(e.buf, d, time.Millisecond, false)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *Entry) Interface(key string, val interface{}) *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, key)
	e.buf = enc.AppendInterface(e.buf, val)
	return e
}

// StackTrace adds stack_trace field to the current context
func (e *Entry) StackTrace() *Entry {
	if e == nil {
		return e
	}
	e.buf = enc.AppendKey(e.buf, "stack_trace")
	e.buf = enc.AppendString(e.buf, getStackTrace())
	return e
}

func handler(e *Entry) {

	for _, h := range e.logger.cacheLeveledHandlers(e.Level) {

		newEntry := copyEntry(e)

		// call hook interface
		for _, hooker := range _logger.hooks {
			_ = hooker(newEntry)
		}

		err := h.BeforeWriting(newEntry)
		if err != nil {
			stdlog.Printf("log: log hook failed: %v", err)
		}

		if len(newEntry.Message) > 0 {
			newEntry.buf = enc.AppendKey(newEntry.buf, "msg")
			newEntry.buf = enc.AppendString(newEntry.buf, newEntry.Message)
		}

		newEntry.buf = enc.AppendEndMarker(newEntry.buf)
		newEntry.buf = enc.AppendLineBreak(newEntry.buf)

		err = h.Write(newEntry.buf)
		if err != nil {
			if ErrorHandler != nil {
				ErrorHandler(err)
			} else {
				stdlog.Printf("log: log write failed: %v", err)
			}
		}

		putEntry(newEntry)
	}

	// hs := e.logger.cacheLeveledHandlers(e.Level)
	// if len(hs) == 0 {
	// 	putEntry(e)
	// 	return
	// }

	// if len(hs) > 1 {
	// 	handlers(e)
	// 	return
	// }

	// h := hs[0]

	// err := h.Hook(e)
	// if err != nil {
	// 	stdlog.Printf("log: log hook failed: %v", err)
	// }

	// if len(e.Message) > 0 {
	// 	e.buf = enc.AppendKey(e.buf, "msg")
	// 	e.buf = enc.AppendString(e.buf, e.Message)
	// }

	// e.buf = enc.AppendEndMarker(e.buf)
	// e.buf = enc.AppendLineBreak(e.buf)

	// err = h.Write(e.buf)
	// if err != nil {
	// 	stdlog.Printf("log: log write failed: %v", err)
	// }

	putEntry(e)
}
