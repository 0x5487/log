package main

import (
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
	//err := errors.New("something bad happened")
	var err error
	log.WithError(err).Error("oops...")
}
