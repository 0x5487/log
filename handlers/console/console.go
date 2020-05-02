package console

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/fatih/color"
	"github.com/jasonsoft/log"
	colorable "github.com/mattn/go-colorable"
)

var colors = [...]*color.Color{
	log.DebugLevel: color.New(color.FgWhite),
	log.InfoLevel:  color.New(color.FgBlue),
	log.WarnLevel:  color.New(color.FgYellow),
	log.ErrorLevel: color.New(color.FgRed),
	log.FatalLevel: color.New(color.FgRed),
	log.PanicLevel: color.New(color.FgRed),
}

var bold = color.New(color.Bold)

// Console is an instance of the console logger
type Console struct {
	mutex  sync.Mutex
	writer io.Writer
}

// New create a new Console instance
func New() *Console {
	return &Console{
		writer: colorable.NewColorableStdout(),
	}
}

// Hook handles the log entry
func (h *Console) Hook(e *log.Entry) error {
	e.Str("level", e.Level.String())

	return nil
}

// Write handles the log entry
func (h *Console) Write(e *log.Entry) error {
	color := colors[e.Level]
	level := e.Level.String()

	kv := map[string]interface{}{}
	err := json.Unmarshal(e.Buffer(), &kv)
	if err != nil {
		return err
	}

	// fmt is not goroutine safe
	// https://stackoverflow.com/questions/14694088/is-it-safe-for-more-than-one-goroutine-to-print-to-stdout
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, _ = color.Fprintf(h.writer, "%s %-50s", bold.Sprintf("%-8s", level), e.Message)

	for k, v := range kv {
		if k == "level" || k == "msg" {
			continue
		}
		fmt.Fprintf(h.writer, " %s=%v", color.Sprint(k), fmt.Sprintf("%v", v))
	}

	fmt.Fprintln(h.writer)

	return nil
}
