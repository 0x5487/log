package console

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jasonsoft/log"
)

const (
	base10 = 10
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

	builder := strings.Builder{}

	if entry.Line > 0 && len(entry.File) > 0 {
		builder.WriteString(fmt.Sprintf("time=\"%s\" level=%s msg=\"%s\" line=\"%d\" file=\"%s\" ", time, level, entry.Message, entry.Line, entry.File))
	} else {
		builder.WriteString(fmt.Sprintf("time=\"%s\" level=%s msg=\"%s\" ", time, level, entry.Message))
	}

	// custom fields to string
	var b []byte
	for key, value := range entry.Fields {
		b = append(b, key...)
		b = append(b, "="...)
		switch t := value.(type) {
		case string:
			b = append(b, t...)
		case int:
			b = strconv.AppendInt(b, int64(t), base10)
		case int8:
			b = strconv.AppendInt(b, int64(t), base10)
		case int16:
			b = strconv.AppendInt(b, int64(t), base10)
		case int32:
			b = strconv.AppendInt(b, int64(t), base10)
		case int64:
			b = strconv.AppendInt(b, t, base10)
		case uint:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint8:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint16:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint32:
			b = strconv.AppendUint(b, uint64(t), base10)
		case uint64:
			b = strconv.AppendUint(b, t, base10)
		case float32:
			b = strconv.AppendFloat(b, float64(t), 'f', -1, 32)
		case float64:
			b = strconv.AppendFloat(b, t, 'f', -1, 64)
		case bool:
			b = strconv.AppendBool(b, t)
		default:
			b = append(b, fmt.Sprintf("%#v", value)...)
		}
		builder.Write(b)
	}

	return builder.String()
}
