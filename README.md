# log
It is a simple log library for golang.  Golang standard context is Supported.

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
	log.SetAppID("TesterApp") // unique id for the app

	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	graylog := gelf.New("tcp://192.168.1.1:12201")
	log.RegisterHandler(graylog, log.AllLevels...)

	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		log.Debug("hello world")
		customFields := log.Fields{
			"city":    "keelung",
			"country": "taiwan",
		}

		log.WithFields(customFields).Info("more info")
		log.Error("oops...")
	}
	duration := int64(time.Since(startTime) / time.Millisecond)
	println(duration)

	time.Sleep(5 * time.Second)
}
```