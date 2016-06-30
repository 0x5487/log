package log

import (
	"sync"
	"time"
)

const (
	keyVal = " %s=%v"
	cutset = "\r\n\t "
)

type CustomFields map[string]interface{}

type Entry struct {
	wg        *sync.WaitGroup
	calldepth int

	Level        Level
	Message      string
	File         string
	Line         int
	Timestamp    time.Time
	CustomFields CustomFields
}

// Consumed lets the Entry and subsequently the Logger
// instance know that it has been used by a handler
func (e *Entry) Consumed() {
	if e.wg != nil {
		e.wg.Done()
	}
}
