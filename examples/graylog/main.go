package main

import (
	"errors"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	"github.com/jasonsoft/log/v2/handlers/gelf"
)

func main() {
	clog := console.New()
	log.AddHandler(clog, log.AllLevels...) // use console handler to log all level log

	graylog := gelf.New("tcp://10.1.1.181:12201")
	log.AddHandler(graylog, log.AllLevels...)

	defer log.Flush()

	logger := log.
		Str("app_id", "santa").
		Str("env", "dev")

	defer log.Trace("time to run").Stop() // use trace to know how long it takes

	log.Debug("hello world1")

	logger = logger.Str("city", "keelung").Str("name", "abc")

	logger.Infof("more info") // log information with custom fileds

	err := errors.New("something bad happened")
	log.Err(err).Error("oops...") // log error struct and print error message
}
