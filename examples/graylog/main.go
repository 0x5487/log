package main

import (
	"errors"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
)

func main() {
	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...) // use console handler to log all level log

	graylog := gelf.New("tcp://10.1.1.181:12201")
	log.RegisterHandler(graylog, log.AllLevels...)
	defer graylog.Flush()

	log.WithDefaultFields(log.Fields{
		"app_id": "santa",
		"host":   "taiwan",
	})

	defer log.Trace("time to run").Stop() // use trace to know how long it takes

	log.Debug("hello world1")

	fields := log.Fields{
		"city":    "keelung",
		"country": "taiwan",
	}
	log.WithFields(fields).Infof("more info") // log information with custom fileds

	err := errors.New("something bad happened")
	log.WithError(err).Error("oops...") // log error struct and print error message
}
