package console

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/fatih/color"
	"github.com/jasonsoft/log"
	colorable "github.com/mattn/go-colorable"
)

var colors = []*color.Color{
	log.DebugLevel: color.New(color.FgWhite),
	log.InfoLevel:  color.New(color.FgBlue),
	log.WarnLevel:  color.New(color.FgYellow),
	log.ErrorLevel: color.New(color.FgRed),
	log.FatalLevel: color.New(color.FgRed),
	log.PanicLevel: color.New(color.FgRed),
}

func levelToColor(level string) *color.Color {
	switch level {
	case "DEBUG":
		return color.New(color.FgWhite)
	case "INFO":
		return color.New(color.FgBlue)
	case "WARN":
		return color.New(color.FgYellow)
	case "ERROR", "PANIC", "FATAL":
		return color.New(color.FgRed)
	default:
		return color.New(color.FgWhite)
	}
}

var bold = color.New(color.Bold)

// Console is an instance of the console logger
type Console struct {
	mutex  sync.Mutex
	writer io.Writer
}

// New create a new Console instance
func New() log.Handler {
	return &Console{
		writer: colorable.NewColorableStdout(),
	}
}

// BeforeWriting handles the log entry
func (h *Console) BeforeWriting(e *log.Entry) error {
	e.Str("level", e.Level.String())

	return nil
}

// Write handles the log entry
func (h *Console) Write(bytes []byte) error {
	kv := map[string]interface{}{}
	err := json.Unmarshal(bytes, &kv)
	if err != nil {
		return err
	}

	level := fmt.Sprintf("%v", kv["level"])
	msg := kv["msg"]
	color := levelToColor(level)

	// sort map by key
	keys := make([]string, 0, len(kv))
	for k := range kv {
		if k == "level" || k == "msg" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// fmt is not goroutine safe
	// https://stackoverflow.com/questions/14694088/is-it-safe-for-more-than-one-goroutine-to-print-to-stdout
	h.mutex.Lock()
	defer h.mutex.Unlock()

	_, _ = color.Fprintf(h.writer, "%s %-50s", bold.Sprintf("%-8s", level), msg)

	for _, k := range keys {
		fmt.Fprintf(h.writer, " %s=%v", color.Sprint(k), fmt.Sprintf("%v", kv[k]))
	}

	fmt.Fprintln(h.writer)

	return nil
}
