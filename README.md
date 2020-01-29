# log
It is a simple sturctured logging package  for Go.

## Features

* easy and configurable
* bulit-in some handlers
* allow to add default fields to every log.  ( ex.  You maybe want to add `app_id` per each app or `env` per each environment)
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
```shell
time="2020-01-29T02:00:53.892Z" level=DEBUG msg="hello world" app_id=santa env=dev
time="2020-01-29T02:00:53.892Z" level=INFO msg="more info" env=dev city=keelung app_id=santa
time="2020-01-29T02:00:53.892Z" level=ERROR msg="oops..." env=dev error=something bad happened app_id=santa
time="2020-01-29T02:00:53.892Z" level=INFO msg="time to run" app_id=santa duration=0s env=dev
```
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

