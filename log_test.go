package log_test

import (
	"context"
	"testing"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"
	"github.com/jasonsoft/log/handlers/gelf"
)

func init() {
	clog := console.New()
	log.RegisterHandler(clog, log.GetLevelsFromMinLevel(log.DebugLevel.String())...)
	graylog := gelf.New("tcp://localhost:12201")
	log.RegisterHandler(graylog, log.GetLevelsFromMinLevel(log.DebugLevel.String())...)
}

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		log.WithFields(log.Fields{"a": 1}).Error("has fields")
		log.WithFields(log.Fields{"a": 2}).Error("has fields 2")
	}
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	ctx = log.NewContext(ctx, log.WithFields(log.Fields{"request_id": 123}))

	logger := log.FromContext(ctx)
	logger.Info("request test")
}

func TestSetLogger(t *testing.T) {
	log.SetLogger(log.WithFields(log.Fields{"global": "globalvar"}))
	log.Info("global test")
}
