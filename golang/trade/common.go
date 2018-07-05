package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

)

const EnableHash bool = true

// Validate chaincode arguments
func checkArguments(args []string, count int) error {
	if len(args) != count {
		return newError("Incorrect number of arguments. Expecting " + strconv.Itoa(count))
	}
	return nil
}

// Macros function
func newError(text string) error {
	fmt.Println(text)
	return errors.New(text)
}




// Hash function
func hash(str string) string {
	if !EnableHash || str == "" {
		return str
	}

	h := sha256.New()
	h.Write([]byte(str))
	s := h.Sum(nil)
	return hex.EncodeToString(s)
}

// Split comma separated string into an array
func splitString(s string) []string {
	var a []string
	if len(s) > 0 {
		a = strings.Split(s, ",")
	}
	return a
}

// Convert string to int64
func stringToInt64(s string) int64 {
	i64, _ := strconv.ParseInt(s, 10, 64)
	return i64
}

// Convert string to float64
func stringToFloat64(s string) float64 {
	f64, _ := strconv.ParseFloat(s, 64)
	return f64
}

// Return the number of milliseconds elapsed since January 1, 1970 UTC
func timeNowInUnixMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Return trne if time is between the specific period
func isTimeBetween(time, from, to int64) bool {
	if from > time {
		return false
	}
	if to == 0 {
		return true
	}
	return time <= to
}

// Return true if strings are the same or arg is empty
func isEmptyOrEqual(arg, s string) bool {
	return arg == "" || arg == s
}

// Return true if the first array is completely contained in the second array.
// There must be at least the same number of duplicate values in second as
// there are in first.
func IsSubset(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value] += 1
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}

