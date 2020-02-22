// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package benchmarks

import (
	"fmt"
	"testing"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/discard"
)

func fakeJasonSoftFields() log.Fields {
	return log.Fields{
		"int":     _tenInts[0],
		"ints":    _tenInts,
		"string":  _tenStrings[0],
		"strings": _tenStrings,
		"time":    _tenTimes[0],
		"times":   _tenTimes,
		"user1":   _oneUser,
		"user2":   _oneUser,
		"users":   _tenUsers,
		"error":   errExample,
	}
}

func BenchmarkSmall(b *testing.B) {
	h := discard.New()
	log.RegisterHandler(h, log.InfoLevel)

	for i := 0; i < b.N; i++ {
		log.Info("login")
	}
}

func BenchmarkMedium(b *testing.B) {
	h := discard.New()
	log.RegisterHandler(h, log.InfoLevel)

	for i := 0; i < b.N; i++ {
		log.WithFields(log.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).Info("upload")
	}
}

func BenchmarkLarge(b *testing.B) {
	h := discard.New()
	log.RegisterHandler(h, log.InfoLevel)

	err := fmt.Errorf("boom")

	for i := 0; i < b.N; i++ {
		log.WithFields(log.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).
			WithFields(log.Fields{
				"some":     "more",
				"data":     "here",
				"whatever": "blah blah",
				"more":     "stuff",
				"context":  "such useful",
				"much":     "fun",
			}).
			WithError(err).Error("upload failed")
	}
}
