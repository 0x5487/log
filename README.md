# log [![GoDoc][doc-img]][doc] [![Build][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![GoReport][report-img]][report]

It is a simple structured logging package for Go. 
## Features
* fast, easy to use, and pretty logging for development
* low to zero allocation
* JSON encoding format
* support multiple output and each output can have their own output level 
* goroutine safety (thread-safe)
* colored text for console handler (`linux`, `mac`, and `windows` are supported)
* `context.Context` integration

## Handlers
* console
* gelf (graylog)
* memory (unit test)
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

	// optional: allow handlers to clear all buffer
	defer log.Flush()

	// use withDefaultFields to add fields to every logs
	log.WithDefaultFields(
		log.Fields{
			"app_id": "santa",
			"env":    "dev",
		},
	)

	// use trace to get how long it takes
	defer log.Trace("time to run").Stop()

	// print message use DEBUG level
	log.Debug("hello world")

	// log information with custom fileds
	log.Str("city", "keelung").Infof("more info")

	// log error struct and print error message
	err := errors.New("something bad happened")
	log.WithError(err).Error("oops...")
}
```
Output

![](colored.png)

## Field Types

### Standard Types

* `Str`
* `Bool`
* `Int`, `Int8`, `Int16`, `Int32`, `Int64`
* `Uint`, `Uint8`, `Uint16`, `Uint32`, `Uint64`
* `Float32`, `Float64`


## Benchmarks

Run on MacBook Pro 15-inch 2018 using go version go1.13.5 windows 10 OS

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


[doc-img]: https://godoc.org/github.com/jasonsoft/log?status.svg
[doc]: https://pkg.go.dev/github.com/jasonsoft/log?tab=doc
[ci-img]: https://github.com/jasonsoft/log/workflows/build/badge.svg
[ci]: https://github.com/jasonsoft/log/actions
[cov-img]: https://codecov.io/gh/jasonsoft/log/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/jasonsoft/log
[report-img]: https://goreportcard.com/badge/github.com/jasonsoft/log
[report]: https://goreportcard.com/report/github.com/jasonsoft/log