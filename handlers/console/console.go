package console

import (
	"fmt"

	"jasonsoft/log"
)

// Console is an instance of the console logger
type Console struct {
}

func New() *Console {
	return &Console{}
}

// Run starts the logger consuming on the returned channed
func (c *Console) Run() chan<- *log.Entry {
	// in a big high traffic app, set a higher buffer
	ch := make(chan *log.Entry, 3000)

	go func(entries <-chan *log.Entry) {
		/*
			for {
				select {
				case message := <-entries:
					println(message.Message)
					message.Consumed()
				default:
					println("no msg")
				}

			}*/

		var e *log.Entry
		for e = range entries {
			msg := FormatFunc(e)
			println(msg)
			e.Consumed()
		}
	}(ch)

	return ch
}

func FormatFunc(entry *log.Entry) string {
	time := entry.Timestamp.String()
	level := entry.Level

	return fmt.Sprintf("%s [%s] %s", time, level, entry.Message)
}
