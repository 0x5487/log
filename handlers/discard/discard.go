// Package discard implements a no-op handler useful for benchmarks and tests.
package discard

import (
	"github.com/jasonsoft/log"
)

// Handler implementation.
type Handler struct{}

// New handler.
func New() *Handler {
	return &Handler{}
}

// Log implements log.Handler.
func (h *Handler) Log(e log.Entry) error {
	return nil
}
