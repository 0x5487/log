package log_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/memory"
	"github.com/stretchr/testify/assert"
)

func TestNoHandler(t *testing.T) {
	log.Info("no handler")
}

func TestPrintf(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	log.Infof("logged in %s", "Tobi")

	e := h.Entries[0]
	assert.Equal(t, "logged in Tobi", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
}

func TestFlush(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.GetLevelsFromMinLevel("debug")...)

	log.Debugf("logged in")

	assert.Equal(t, 1, len(h.Entries))

	log.Flush()
	assert.Equal(t, 0, len(h.Entries))
}

func TestLevels(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel)

	log.Debug("debug1")
	log.Info("info1")
	log.Warn("warn1")
	log.Error("error1")

	assert.Equal(t, 4, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "debug1", e.Message)
	assert.Equal(t, log.DebugLevel, e.Level)

	e = h.Entries[1]
	assert.Equal(t, "info1", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)

	e = h.Entries[2]
	assert.Equal(t, "warn1", e.Message)
	assert.Equal(t, log.WarnLevel, e.Level)

	e = h.Entries[3]
	assert.Equal(t, "error1", e.Message)
	assert.Equal(t, log.ErrorLevel, e.Level)
}

func TestContext(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	ctx := context.Background()
	ctx = log.NewContext(ctx, log.WithFields(log.Fields{"request_id": 123}))

	logger := log.FromContext(ctx)
	logger.Warnf("request test")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "request test", e.Message)
	assert.Equal(t, log.WarnLevel, e.Level)
	assert.Equal(t, log.Fields{"request_id": 123}, e.Fields)
}

func TestWithFields(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	logger := log.WithFields(log.Fields{"file": "sloth.png"})
	logger.Info("upload complete")
	logger1 := logger.WithField("source", "machine1")
	logger1.Debug("uploading")
	logger.Warn("warning la")

	assert.Equal(t, 3, len(h.Entries))

	// info
	e := h.Entries[0]
	assert.Equal(t, "upload complete", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, log.Fields{"file": "sloth.png"}, e.Fields)

	// debug
	e = h.Entries[1]
	assert.Equal(t, "uploading", e.Message)
	assert.Equal(t, log.DebugLevel, e.Level)
	assert.Equal(t, log.Fields{"file": "sloth.png", "source": "machine1"}, e.Fields)

	// warn
	e = h.Entries[2]
	assert.Equal(t, "warning la", e.Message)
	assert.Equal(t, log.WarnLevel, e.Level)
	assert.Equal(t, log.Fields{"file": "sloth.png"}, e.Fields)
}

func TestWithField(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	logger := log.WithField("file", "sloth.png").WithField("user", "Tobi")
	logger.Debug("uploading")
	logger.Info("upload complete")

	assert.Equal(t, 2, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "uploading", e.Message)
	assert.Equal(t, log.DebugLevel, e.Level)
	assert.Equal(t, log.Fields{"file": "sloth.png", "user": "Tobi"}, e.Fields)
}

func TestWithError(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	err := errors.New("something bad happened")
	log.WithError(err).Errorf("too bad %s", err.Error())

	err = nil
	log.WithError(err).Error("err is nil")

	assert.Equal(t, 2, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "too bad something bad happened", e.Message)
	assert.Equal(t, log.ErrorLevel, e.Level)
	assert.Equal(t, log.Fields{"error": "something bad happened"}, e.Fields)

	e = h.Entries[1]
	assert.Equal(t, "err is nil", e.Message)
	assert.Equal(t, log.ErrorLevel, e.Level)
	assert.Equal(t, log.Fields{}, e.Fields)
}

func TestTrace(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	func() (err error) {
		defer log.Trace("upload").Stop()
		return nil
	}()

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "upload", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.IsType(t, e.Fields["duration"], "0")

}

func TestWithDefaultFields(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	log.WithDefaultFields(log.Fields{
		"app_id": "santa",
		"env":    "dev",
	})
	log.Info("upload complete")

	logger := log.WithFields(log.Fields{"file": "sloth.png"})
	logger.Debugf("debug test")

	assert.Equal(t, 2, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "upload complete", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, e.Fields, log.Fields{"app_id": "santa", "env": "dev"})

	e = h.Entries[1]
	assert.Equal(t, "debug test", e.Message)
	assert.Equal(t, log.DebugLevel, e.Level)
	assert.Equal(t, e.Fields, log.Fields{"app_id": "santa", "env": "dev", "file": "sloth.png"})

}
