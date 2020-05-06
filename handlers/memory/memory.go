// Package memory implements an in-memory handler useful for testing, as the
// entries can be accessed after writes.
package memory

import (
	"sync"

	"github.com/jasonsoft/log"
)

// Handler implementation.
type Handler struct {
	mu  sync.Mutex
	Out []byte
}

// New handler.
func New() *Handler {
	return &Handler{
		Out: make([]byte, 500),
	}
}

// BeforeWriting implements log.Handler.
func (h *Handler) BeforeWriting(e *log.Entry) error {
	e.Str("level", e.Level.String())

	return nil
}

// Write implements log.Handler.
func (h *Handler) Write(bytes []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Out = bytes
	return nil
}

// Flush clear all buffer
func (h *Handler) Flush() error {
	h.Out = []byte{}
	return nil
}
