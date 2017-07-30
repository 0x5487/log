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
	time := entry.Timestamp.Format("2006-01-02T15:04:05.999Z")
	level := entry.Level.String()

	strFields := ""
	for key, value := range entry.Fields {

		switch value.(type) {
		case string:
			strFields += value.(string)
		default:
			strFields += fmt.Sprintf("%#v ", value)
		}

		strFields += fmt.Sprintf("%s=%s ", key, value)
	}

	var result string
	if entry.Line > 0 && len(entry.File) > 0 {
		result = fmt.Sprintf("time=\"%s\" level=%s msg=\"%s\" line=\"%d\" file=\"%s\"  ", time, level, entry.Message, entry.Line, entry.File)

	} else {
		result = fmt.Sprintf("time=\"%s\" level=%s msg=\"%s\" ", time, level, entry.Message)

	}

	if len(strFields) > 0 {
		result += strings.TrimSpace(strFields)
	}

	return result
}
