package log

import (
	"context"
	"errors"
	"testing"
)

func init() {

}

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := errors.New("something bad...")
		Error(err)
	}
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	ctx = NewContext(ctx, WithFields(Fields{"request_id": 123}))

	logger := FromContext(ctx)
	logger.Info("request test")
}
