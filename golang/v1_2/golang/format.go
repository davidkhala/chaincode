package golang

import (
	"strconv"
	"encoding/json"
	"errors"
	"time"
)

const (
	// Seconds field of the earliest valid Timestamp.
	// This is time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	minValidSeconds = -62135596800
	// Seconds field just after the latest valid Timestamp.
	// This is time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	maxValidSeconds = 253402300800
)

func ToInt(bytes []byte) int {
	if bytes == nil {
		return 0
	}
	i, err := strconv.Atoi(string(bytes))
	PanicError(err)
	return i
}
func ToString(integer interface{}) string {
	return strconv.FormatInt(integer.(int64), 10)
}

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
func PanicString(err string) {
	if err != "" {
		panic(errors.New(err))
	}
}
func UnixMilliSecond(t time.Time) TimeLong {
	return TimeLong(t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
}
type TimeLong int64

/**
	a wrapper to panic Unmarshal(non-pointer v)
 */
func FromJson(jsonString []byte, v interface{}) {
	err := json.Unmarshal(jsonString, v)
	PanicError(err);
}

func ToJson(v interface{}) []byte {
	result, err := json.Marshal(v)
	PanicError(err)
	return result
}
