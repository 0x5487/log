# log
It is a simple log library for golang

## Example

```golang
package main

import (
	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
)

func main() {
	cLog := console.New()

	logger := log.New()
	logger.RegisterHandler(cLog, log.DebugLevel, log.InfoLevel, log.ErrorLevel)

	logger.Debug("hello world")

	customFields := log.Fields{
		"car":     "bmw",
		"country": "taiwan",
	}

	logger.WithFields(customFields).Info("more info")

	logger.Error("oops...")
}
```