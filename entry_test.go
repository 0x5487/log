package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
	Age  int
}

func TestEntryFields(t *testing.T) {
	entry := newEntry(_logger, nil)

	time1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	time2, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+08:00")

	entry = entry.
		Str("hello", "world").
		Strs("strs", []string{"str1", "str2"}).
		Bool("is_enabled", true).
		Int("int", 1).
		Ints("ints", []int{1, 2}).
		Int8("int8", int8(2)).
		Int16("int16", int16(3)).
		Int32("int32", int32(4)).
		Int64("int64", int64(5)).
		Uint("uint", uint(6)).
		Uint8("uint8", uint8(7)).
		Uint16("uint16", uint16(8)).
		Uint32("uint32", uint32(9)).
		Uint64("uint64", uint64(10)).
		Float32("float32", float32(11.123)).
		Float64("float64", float64(12.123)).
		Time("time", time1).
		Times("times", []time.Time{time1, time2}).
		Interface("person", Person{})

	entry.Debug("debug")

	// t.Log(string(entry.buf))
	assert.Equal(t, DebugLevel, entry.Level)
	assert.Equal(t, "debug", entry.Message)
	assert.Equal(t, `{"hello":"world","strs":["str1","str2"],"is_enabled":true,"int":1,"ints":[1,2],"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.123,"float64":12.123,"time":"2012-11-01T22:08:41Z","times":["2012-11-01T22:08:41Z","2012-11-01T22:08:41+08:00"],"person":{"Name":"","Age":0}`, string(entry.buf))

}
