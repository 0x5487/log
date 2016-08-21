# log
It is a simple log library for golang

## Supported handlers
* console
* gelf (graylog)

## Example

```go
package main

import (
	"time"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
)

func main() {
	clog := console.New()
	graylog := gelf.New("tcp://192.168.1.1:12201")

	logger := log.New()
	logger.SetAppID("TesterApp") // unique id for the app
	logger.RegisterHandler(clog, log.AllLevels...)
	logger.RegisterHandler(graylog, log.AllLevels...)

	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		logger.Debug("hello world")
		customFields := log.Fields{
			"city":     "keelung",
			"country": "taiwan",
		}

		logger.WithFields(customFields).Info("more info")
		logger.Error("oops...")
	}
	duration := int64(time.Since(startTime) / time.Millisecond)
	println(duration)
}

```