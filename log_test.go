package log_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/discard"
	"github.com/jasonsoft/log/handlers/memory"
	"github.com/tj/assert"
)

func TestNoHandler(t *testing.T) {
	log.Info("no handler 1")
	log.Warnf("no handler 2")
}

func TestRegisterHandlers(t *testing.T) {
	h1 := memory.New()
	log.RegisterHandler(h1, log.AllLevels...)

	h2 := memory.New()
	log.RegisterHandler(h2, log.AllLevels...)

	log.Info("info")
	assert.Equal(t, `{"level":"INFO","msg":"info"}`+"\n", string(h1.Out))
	assert.Equal(t, `{"level":"INFO","msg":"info"}`+"\n", string(h2.Out))
}

func TestLog(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	log.Debug("debug")
	assert.Equal(t, `{"level":"DEBUG","msg":"debug"}`+"\n", string(h.Out))

	log.Debugf("debug %s", "hello")
	assert.Equal(t, `{"level":"DEBUG","msg":"debug hello"}`+"\n", string(h.Out))

	log.Info("info")
	assert.Equal(t, `{"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Infof("info %s", "hello")
	assert.Equal(t, `{"level":"INFO","msg":"info hello"}`+"\n", string(h.Out))

	log.Warn("warn")
	assert.Equal(t, `{"level":"WARN","msg":"warn"}`+"\n", string(h.Out))

	log.Warnf("warn %s", "hello")
	assert.Equal(t, `{"level":"WARN","msg":"warn hello"}`+"\n", string(h.Out))

	log.Error("error")
	assert.Equal(t, `{"level":"ERROR","msg":"error"}`+"\n", string(h.Out))

	log.Errorf("error %s", "hello")
	assert.Equal(t, `{"level":"ERROR","msg":"error hello"}`+"\n", string(h.Out))

	t.Run("test panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if ok {
					assert.Equal(t, "panic", err.Error())
				}
			}
			assert.Equal(t, `{"level":"PANIC","msg":"panic"}`+"\n", string(h.Out))
		}()
		log.Panic("panic")
	})

	t.Run("test panicf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if ok {
					assert.Equal(t, "panic hello", err.Error())
				}
			}
			assert.Equal(t, `{"level":"PANIC","msg":"panic hello"}`+"\n", string(h.Out))
		}()
		log.Panicf("panic %s", "hello")
	})

}

func TestContext(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	logger := log.Str("app", "stant")

	logger.Str("a", "b").Info("hello world")
	//t.Log(string(h.Out))
	assert.Equal(t, `{"app":"stant","a":"b","level":"INFO","msg":"hello world"}`+"\n", string(h.Out))

	logger.Str("c", "d").Info("hello world")
	assert.Equal(t, `{"app":"stant","c":"d","level":"INFO","msg":"hello world"}`+"\n", string(h.Out))
}

func TestFlush(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.GetLevelsFromMinLevel("debug")...)

	log.Debug("flush")
	log.Flush()
	assert.Equal(t, 0, len(h.Out))
}

func TestLevels(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel)

	log.Debug("debug1")
	//t.Log(string(h.Out))
	assert.Equal(t, `{"level":"DEBUG","msg":"debug1"}`+"\n", string(h.Out))

	log.Info("info1")
	assert.Equal(t, `{"level":"INFO","msg":"info1"}`+"\n", string(h.Out))

	log.Warn("warn1")
	assert.Equal(t, `{"level":"WARN","msg":"warn1"}`+"\n", string(h.Out))

	log.Error("error1")
	assert.Equal(t, `{"level":"ERROR","msg":"error1"}`+"\n", string(h.Out))

}

func TestStdContext(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	t.Run("create new context", func(t *testing.T) {
		ctx := context.Background()
		ctx = log.Str("request_id", "abc").WithContext(ctx)

		logger := log.FromContext(ctx)
		logger.Debug("test")
		assert.Equal(t, `{"request_id":"abc","level":"DEBUG","msg":"test"}`+"\n", string(h.Out))

		logger.Str("app", "santa").Debugf("debug %s", "hello")
		//t.Log(string(h.Out))
		assert.Equal(t, `{"request_id":"abc","app":"santa","level":"DEBUG","msg":"debug hello"}`+"\n", string(h.Out))
	})

	t.Run("create blank context", func(t *testing.T) {
		ctx := context.Background()
		logger := log.FromContext(ctx)
		logger.Info("test")
		//t.Log(string(h.Out))
		assert.Equal(t, `{"level":"INFO","msg":"test"}`+"\n", string(h.Out))

	})
}

func TestStandardFields(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.GetLevelsFromMinLevel("debug")...)

	logger := log.
		Str("hello", "world").
		Bool("is_enabled", true).
		Int("int", 1).
		Int8("int8", int8(2)).
		Int16("int16", int16(3)).
		Int32("int32", int32(4)).
		Int64("int64", int64(5)).
		Uint("uint", uint(6)).
		Uint8("uint8", uint8(7)).
		Uint16("uint16", uint16(8)).
		Uint32("uint32", uint32(9)).
		Uint64("uint64", uint64(10)).
		Float32("float32", float32(11.123)).
		Float64("float64", float64(12.123))

	logger.Debug("debug")

	//t.Log(string(h.Out))
	assert.Equal(t, `{"hello":"world","is_enabled":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.123,"float64":12.123,"level":"DEBUG","msg":"debug"}`+"\n", string(h.Out))

}

func TestAdvancedFields(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	err := errors.New("something bad happened")
	log.Err(err).Error("too bad")

	//t.Log(string(h.Out))
	assert.Equal(t, `{"error":"something bad happened","level":"ERROR","msg":"too bad"}`+"\n", string(h.Out))

}

// func TestTrace(t *testing.T) {
// 	h := memory.New()
// 	log.RegisterHandler(h, log.AllLevels...)

// 	func() (err error) {
// 		defer log.Trace("trace").Stop()
// 		return nil
// 	}()

// 	t.Log("aa" + string(h.Out))

// }

// func TestWithDefaultFields(t *testing.T) {
// 	h := memory.New()
// 	log.RegisterHandler(h, log.AllLevels...)

// 	log.WithDefaultFields(log.Fields{
// 		"app_id": "santa",
// 		"env":    "dev",
// 	})
// 	log.Info("upload complete")

// 	logger := log.WithFields(log.Fields{"file": "sloth.png"})
// 	logger.Debugf("debug test")

// 	assert.Equal(t, 2, len(h.Entries))

// 	e := h.Entries[0]
// 	assert.Equal(t, "upload complete", e.Message)
// 	assert.Equal(t, log.InfoLevel, e.Level)
// 	assert.Equal(t, e.Fields, log.Fields{"app_id": "santa", "env": "dev"})

// 	e = h.Entries[1]
// 	assert.Equal(t, "debug test", e.Message)
// 	assert.Equal(t, log.DebugLevel, e.Level)
// 	assert.Equal(t, e.Fields, log.Fields{"app_id": "santa", "env": "dev", "file": "sloth.png"})

// }

func TestGoroutineSafe(t *testing.T) {

	logger := log.Str("request_id", "abc")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		logger.Str("name", "abc")
	}()
	go func() {
		defer wg.Done()
		logger.Str("name", "xyz")
	}()
	wg.Wait()
}

func BenchmarkDisabledAddingFields(b *testing.B) {
	b.Logf("Logging without any structured context.")

	b.Run("jasnosoft/log", func(b *testing.B) {
		h := discard.New()
		log.RegisterHandler(h, log.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Info("hello world")
			}
		})
	})

}
