package log_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/discard"
	"github.com/jasonsoft/log/v2/handlers/memory"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
	Age  int
}

func TestNoHandler(t *testing.T) {
	log.Info("no handler 1")
	log.Warnf("no handler 2")
}

func TestAddHandlers(t *testing.T) {
	log.RemoveAllHandlers()
	h1 := memory.New()
	log.AddHandler(h1, log.AllLevels...)

	h2 := memory.New()
	log.AddHandler(h2, log.AllLevels...)

	log.Info("info")
	assert.Equal(t, `{"level":"INFO","msg":"info"}`+"\n", string(h1.Out))
	assert.Equal(t, `{"level":"INFO","msg":"info"}`+"\n", string(h2.Out))
}

func TestLog(t *testing.T) {
	log.RemoveAllHandlers()
	h := memory.New()
	log.AddHandler(h, log.AllLevels...)

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
	log.RemoveAllHandlers()
	h := memory.New()
	log.AddHandler(h, log.AllLevels...)

	logger := log.Str("app", "stant")

	logger.Str("a", "b").Info("hello world")
	//t.Log(string(h.Out))
	assert.Equal(t, `{"app":"stant","a":"b","level":"INFO","msg":"hello world"}`+"\n", string(h.Out))

	logger.Bool("bool", true).Info("hello world")
	assert.Equal(t, `{"app":"stant","bool":true,"level":"INFO","msg":"hello world"}`+"\n", string(h.Out))

	log.Int("int", 1).Info("info")
	assert.Equal(t, `{"int":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Int8("int8", 1).Info("info")
	assert.Equal(t, `{"int8":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Int16("int16", 1).Info("info")
	assert.Equal(t, `{"int16":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Int32("int32", 1).Info("info")
	assert.Equal(t, `{"int32":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Int64("int64", 1).Info("info")
	assert.Equal(t, `{"int64":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Uint("uint", 1).Info("info")
	assert.Equal(t, `{"uint":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Uint8("uint8", 1).Info("info")
	assert.Equal(t, `{"uint8":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Uint16("uint16", 1).Info("info")
	assert.Equal(t, `{"uint16":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Uint32("uint32", 1).Info("info")
	assert.Equal(t, `{"uint32":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Uint64("uint64", 1).Info("info")
	assert.Equal(t, `{"uint64":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Float32("float32", 1).Info("info")
	assert.Equal(t, `{"float32":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

	log.Float64("float64", 1).Info("info")
	assert.Equal(t, `{"float64":1,"level":"INFO","msg":"info"}`+"\n", string(h.Out))

}

func TestFlush(t *testing.T) {
	log.RemoveAllHandlers()
	h := memory.New()
	log.AddHandler(h, log.GetLevelsFromMinLevel("debug")...)

	log.Debug("flush")
	log.Flush()
	assert.Equal(t, 0, len(h.Out))
}

func TestLevels(t *testing.T) {
	log.RemoveAllHandlers()

	levels := log.GetLevelsFromMinLevel("debug")
	assert.Equal(t, log.AllLevels, levels)

	levels = log.GetLevelsFromMinLevel("")
	assert.Equal(t, log.AllLevels, levels)

	levels = log.GetLevelsFromMinLevel("info")
	assert.Equal(t, []log.Level{log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.PanicLevel, log.FatalLevel}, levels)

	levels = log.GetLevelsFromMinLevel("warn")
	assert.Equal(t, []log.Level{log.WarnLevel, log.ErrorLevel, log.PanicLevel, log.FatalLevel}, levels)

	levels = log.GetLevelsFromMinLevel("error")
	assert.Equal(t, []log.Level{log.ErrorLevel, log.PanicLevel, log.FatalLevel}, levels)

	levels = log.GetLevelsFromMinLevel("panic")
	assert.Equal(t, []log.Level{log.PanicLevel, log.FatalLevel}, levels)

	levels = log.GetLevelsFromMinLevel("fatal")
	assert.Equal(t, []log.Level{log.FatalLevel}, levels)
}

func TestStdContext(t *testing.T) {
	log.RemoveAllHandlers()

	h := memory.New()
	log.AddHandler(h, log.AllLevels...)

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
	log.RemoveAllHandlers()

	h := memory.New()
	log.AddHandler(h, log.GetLevelsFromMinLevel("debug")...)

	time1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	time2, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+08:00")

	logger := log.
		Str("hello", "world").
		Strs("strs", []string{"str1", "str2"}).
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
		Float64("float64", float64(12.123)).
		Time("time", time1).
		Times("times", []time.Time{time1, time2}).
		Interface("person", Person{})

	logger.Debug("debug")

	//t.Log(string(h.Out))
	assert.Equal(t, `{"hello":"world","strs":["str1","str2"],"is_enabled":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.123,"float64":12.123,"time":"2012-11-01T22:08:41Z","times":["2012-11-01T22:08:41Z","2012-11-01T22:08:41+08:00"],"person":{"Name":"","Age":0},"level":"DEBUG","msg":"debug"}`+"\n", string(h.Out))

}

func TestAdvancedFields(t *testing.T) {
	log.RemoveAllHandlers()

	h := memory.New()
	log.AddHandler(h, log.AllLevels...)

	err := errors.New("something bad happened")
	log.Err(err).Error("too bad")

	//t.Log(string(h.Out))
	assert.Equal(t, `{"error":"something bad happened","level":"ERROR","msg":"too bad"}`+"\n", string(h.Out))

}

// func TestTrace(t *testing.T) {
// 	h := memory.New()
// 	log.AddHandler(h, log.AllLevels...)

// 	func() (err error) {
// 		defer log.Trace("trace").Stop()
// 		return nil
// 	}()
// 	assert.Equal(t, `{"duration":0,"level":"INFO","msg":"trace"}`+"\n", string(h.Out))
// }

type AppHook struct {
}

func (h *AppHook) Hook(e *log.Entry) error {
	e.Str("app_id", "santa").Str("env", "dev")
	return nil
}

func TestHook(t *testing.T) {
	log.RemoveAllHandlers()

	h := memory.New()
	log.AddHandler(h, log.AllLevels...)

	log.AddHook(func(e *log.Entry) error {
		e.Str("app_id", "santa").Str("env", "dev")
		return nil
	})

	log.Info("upload complete")

	//t.Log(string(h.Out))
	assert.Equal(t, `{"app_id":"santa","env":"dev","level":"INFO","msg":"upload complete"}`+"\n", string(h.Out))
}

func TestGoroutineSafe(t *testing.T) {
	log.RemoveAllHandlers()

	h := memory.New()
	log.AddHandler(h, log.AllLevels...)

	logger := log.Str("request_id", "abc")

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_ = logger.Str("name", "abc")
		logger.Info("test")
	}()
	go func() {
		defer wg.Done()
		_ = logger.Str("name", "xyz")
		logger.Info("test")
	}()
	wg.Wait()
}

func TestSaveToDefaultContext(t *testing.T) {
	log.RemoveAllHandlers()
	h := memory.New()
	log.AddHandler(h, log.AllLevels...)

	log.
		Str("app", "santa").
		Str("env", "dev").
		SaveToDefault()

	log.Debug("hello")
	assert.Equal(t, `{"app":"santa","env":"dev","level":"DEBUG","msg":"hello"}`+"\n", string(h.Out))

	log.Bool("answer", true).SaveToDefault()
	log.Int32("count", 3).Info("hello2")

	//t.Log(string(h.Out))
	assert.Equal(t, `{"app":"santa","env":"dev","answer":true,"count":3,"level":"INFO","msg":"hello2"}`+"\n", string(h.Out))
}

func BenchmarkDisabledAddingFields(b *testing.B) {
	b.Logf("Logging without any structured context.")

	b.Run("jasnosoft/log", func(b *testing.B) {
		h := discard.New()
		log.AddHandler(h, log.InfoLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				log.Info("hello world")
			}
		})
	})
}
