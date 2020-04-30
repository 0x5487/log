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
