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

## Example

```go
package main

import (
	"errors"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
)

func main() {
	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...) // use console handler to log all level log

	defer log.Trace("time to run").Stop() // use trace to know how long it takes

	log.Debug("hello world")

	fields := log.Fields{
		"city":    "keelung",
		"country": "taiwan",
	}
	log.WithFields(fields).Infof("more info") // log information with custom fileds

	err := errors.New("something bad happened")
	log.WithError(err).Error("oops...") // log error struct and print error message
}
```
Output
```shell
time="2020-01-26T03:48:58.359Z" level=DEBUG msg="hello world"
time="2020-01-26T03:48:58.359Z" level=INFO msg="more info" city=keelung country=taiwan
time="2020-01-26T03:48:58.359Z" level=ERROR msg="oops..." error=something bad happened
time="2020-01-26T03:48:58.36Z" level=INFO msg="time to run" duration=983.1µs
```
## Benchmarks

```shell
go test -bench=. -benchmem -run=^bb -v
goos: windows
goarch: amd64
pkg: github.com/jasonsoft/log
BenchmarkSmall-12       13276719                95.4 ns/op             0 B/op          0 allocs/op
BenchmarkMedium-12       1000000              1001 ns/op             336 B/op          2 allocs/op
BenchmarkLarge-12         203383              5320 ns/op            2183 B/op          9 allocs/op
PASS
ok      github.com/jasonsoft/log        4.399s
```

