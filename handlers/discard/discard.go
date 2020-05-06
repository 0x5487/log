// Package discard implements a no-op handler useful for benchmarks and tests.
package discard

import (
	"github.com/jasonsoft/log"
)

// Handler implementation.
type Handler struct{}

// New handler.
func New() log.Handler {
	return &Handler{}
}

// BeforeWriting implements log.Handler.
func (h *Handler) BeforeWriting(e *log.Entry) error {
	e.Str("level", e.Level.String())

	return nil
}

// Log implements log.Handler.
func (h *Handler) Write(bytes []byte) error {
	return nil
}
