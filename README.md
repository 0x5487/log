# log
It is a simple log library for golang.  Golang standard context is Supported.

## Supported handlers
* console
* gelf (graylog)
* memory (unit test purpose)
* discard (benchmark)

## Example

```go
package main

import (
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
)

func main() {
	log.SetAppID("TesterApp") // unique id for the app

	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	// send log to graylog server
	// graylog := gelf.New("tcp://192.168.1.1:12201")
	// log.RegisterHandler(graylog, log.AllLevels...)

	defer log.Trace("time to run").Stop()
	for i := 0; i < 2; i++ {
		log.Debug("hello world")
		customFields := log.Fields{
			"city":    "keelung",
			"country": "taiwan",
		}

		log.WithFields(customFields).Info("more info")
		log.Error("oops...")
	}
}
```

## How to run benchmark
```
go test --bench=. -benchmem
```

