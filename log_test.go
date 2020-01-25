package log_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/discard"
	"github.com/jasonsoft/log/handlers/memory"
	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	log.Infof("logged in %s", "Tobi")

	e := h.Entries[0]
	assert.Equal(t, "logged in Tobi", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
}

func TestLevels(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.InfoLevel)

	log.Debug("uploading")
	log.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "upload complete", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
}

func TestContext(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	ctx := context.Background()
	ctx = log.NewContext(ctx, log.WithFields(log.Fields{"request_id": 123}))

	logger := log.FromContext(ctx)
	logger.Info("request test")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "request test", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, log.Fields{"request_id": 123}, e.Fields)
}

func TestWithFields(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	logger := log.WithFields(log.Fields{"file": "sloth.png"})
	logger.Info("upload complete")
	logger = logger.WithField("source", "machine1")
	logger.Debug("uploading")

	assert.Equal(t, 2, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "upload complete", e.Message)
	assert.Equal(t, log.InfoLevel, e.Level)
	assert.Equal(t, e.Fields, log.Fields{"file": "sloth.png"})

	e = h.Entries[1]
	assert.Equal(t, "uploading", e.Message)
	assert.Equal(t, log.DebugLevel, e.Level)
	assert.Equal(t, e.Fields, log.Fields{"file": "sloth.png", "source": "machine1"})
}

func TestWithError(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	err := errors.New("something bad happened")
	log.WithError(err).Error("too bad")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, "too bad", e.Message)
	assert.Equal(t, log.ErrorLevel, e.Level)
	assert.Equal(t, log.Fields{"error": "something bad happened"}, e.Fields)
}

func TestTrace(t *testing.T) {
	h := memory.New()
	log.RegisterHandler(h, log.AllLevels...)

	func() (err error) {
		defer log.WithField("file", "sloth.png").Trace("upload").Stop()
		return nil
	}()

	assert.Equal(t, 1, len(h.Entries))
	{
		e := h.Entries[0]
		assert.Equal(t, "upload", e.Message)
		assert.Equal(t, log.InfoLevel, e.Level)
		assert.Equal(t, e.Fields["file"], "sloth.png")
		assert.IsType(t, e.Fields["duration"], "0")
	}
}

func BenchmarkSmall(b *testing.B) {
	h := discard.New()
	log.RegisterHandler(h, log.InfoLevel)

	for i := 0; i < b.N; i++ {
		log.Info("login")
	}
}

func BenchmarkMedium(b *testing.B) {
	h := discard.New()
	log.RegisterHandler(h, log.InfoLevel)

	for i := 0; i < b.N; i++ {
		log.WithFields(log.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).Info("upload")
	}
}

func BenchmarkLarge(b *testing.B) {
	h := discard.New()
	log.RegisterHandler(h, log.InfoLevel)

	err := fmt.Errorf("boom")

	for i := 0; i < b.N; i++ {
		log.WithFields(log.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).
			WithFields(log.Fields{
				"some":     "more",
				"data":     "here",
				"whatever": "blah blah",
				"more":     "stuff",
				"context":  "such useful",
				"much":     "fun",
			}).
			WithError(err).Error("upload failed")
	}
}
