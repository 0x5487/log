package main

import (
	"errors"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
)

func main() {
	// use console handler to log all level logs
	clog := console.New()
	log.AddHandler(clog, log.AllLevels...)

	// optional: allow handlers to clear all buffer
	defer log.Flush()

	// use trace to get how long it takes
	defer log.Trace("time to run").Stop()

	logger := log.
		Str("app_id", "santa").
		Str("env", "dev")

	// print message use DEBUG level
	logger.Debug("hello world")

	// log information with custom fileds
	logger.Str("city", "keelung").Info("more info")

	// log error struct and print error message
	err := errors.New("something bad happened")
	logger.Err(err).Error("oops...")
}
