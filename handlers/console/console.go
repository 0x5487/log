package console

import (
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

// Log handles the log entry
func (h *Console) Log(e log.Entry) error {
	color := colors[e.Level]
	level := e.Level.String()
	names := e.Fields.Names()

	// fmt is not goroutine safe
	// https://stackoverflow.com/questions/14694088/is-it-safe-for-more-than-one-goroutine-to-print-to-stdout
	h.mutex.Lock()
	defer h.mutex.Unlock()

	color.Fprintf(h.writer, "%s %-50s", bold.Sprintf("%-8s", level), e.Message)

	for _, name := range names {
		fmt.Fprintf(h.writer, " %s=%v", color.Sprint(name), e.Fields.Get(name))
	}

	fmt.Fprintln(h.writer)

	return nil
}
