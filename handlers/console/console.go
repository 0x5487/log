package console

import (
	"fmt"
	"strings"

	"github.com/jasonsoft/log"
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
		for {
			select {
			case e := <-entries:
				msg := FormatFunc(e)
				println(msg)
				e.Consumed()
			default:
			}
		}
	}(ch)

	return ch
}

func FormatFunc(entry *log.Entry) string {
	time := entry.Timestamp.Format("2006-01-02 15:04:05.999")
	level := entry.Level.String()

	strFields := ""
	for key, value := range entry.Fields {
		strFields += key + "=" + value.(string) + " "
	}

	result := fmt.Sprintf("time=\"%s\" level=%s msg=\"%s\" ", time, level, entry.Message)

	if len(strFields) > 0 {
		result += strings.TrimSpace(strFields)
	}

	return result
}
