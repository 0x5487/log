# log
It is a simple sturctured logging package  for Go.

## Features

* easy and configurable
* bulit-in some handlers
* allow to add default fields to every log.  ( ex.  You maybe want to add `app_id` per each app or `env` per each environment)
* colored text for console handler
* trace duration
* work with error interface 
* golang standard context is supported

## handlers
* console
* gelf (graylog)
* memory (unit test purpose)
* discard (benchmark)

## Installation
Use go get 

```go
go get -u github.com/jasonsoft/log
```

## Example

```go
package main

import (
	"errors"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
)

func main() {
	// use console handler to log all level logs
	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	// use withDefaultFields to add fields to every logs
	log.WithDefaultFields(
		log.Fields{
			"app_id": "santa",
			"env":    "dev",
		},
	)

	// use trace to get how long it takes
	defer log.Trace("time to run").Stop()
	log.Debug("hello world")

	// log information with custom fileds
	fields := log.Fields{
		"city": "keelung",
	}
	log.WithFields(fields).Infof("more info")

	// log error struct and print error message
	err := errors.New("something bad happened")
	log.WithError(err).Error("oops...")
}
```
Output

![](colored.png)



## Benchmarks

Run on Macbook Pro 15-inch 2018 using go version go1.13.5 windows 10 os

```shell
go test -bench=. -benchmem -run=^bb -v
goos: windows
goarch: amd64
pkg: github.com/jasonsoft/log
BenchmarkSmall-12       13483690                82.6 ns/op             0 B/op          0 allocs/op
BenchmarkMedium-12       2489635               605 ns/op             336 B/op          2 allocs/op
BenchmarkLarge-12         479955              2802 ns/op            2183 B/op          9 allocs/op
PASS
ok      github.com/jasonsoft/log        4.604s
```

