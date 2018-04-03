package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const EnableHash bool = true


// Get record count from state database

// Validate chaincode arguments
func checkArguments(args []string, count int) error {
	if len(args) != count {
		return newError("Incorrect number of arguments. Expecting " + strconv.Itoa(count))
	}
	return nil
}

// Create an array of partial composite keys
func createPartialCompositeKey(args []string) []string {
	var keys []string
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			break
		}

		keys = append(keys, args[i])
	}
	return keys
}

// Macros function
func accAllowed(r interface{}) error {
	fmt.Println("Allowed")
	emptyOmittedFields(r)
	return nil
}

// Macros function
func newError(text string) error {
	fmt.Println(text)
	return errors.New(text)
}

// Ignore json field if "omitempty" is set
func emptyOmittedFields(v interface{}) interface{} {
	e := reflect.ValueOf(v).Elem()
	for i := 0; i < e.NumField(); i++ {
		json := e.Type().Field(i).Tag.Get("json")
		if strings.Contains(json, "omitempty") {
			f := e.Field(i)
			if f.IsValid() {
				f.SetString("")
			}
		}
	}
	return e.Interface()
}

// Hash function
func Hash(str string) string {
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
func isSubset(first, second []string) bool {
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

// =============================================================================
// Read - read a generic variable from ledger
//
// Inputs - Array of strings
//    0
//   key
//
// Returns - []byte
// =============================================================================
func read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting read")

	err := checkArguments(args, 1)
	if err != nil {
		return shim.Error(err.Error())
	}

	// read the variable from ledger
	key := args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}

	// send it onward
	fmt.Println("- end read")
	return shim.Success(valAsbytes)
}

// =============================================================================
// write() - genric write variable into ledger
//
// Inputs - Array of strings
//    0   ,    1
//   key  ,  value
// =============================================================================
func write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("starting write")

	err := checkArguments(args, 2)
	if err != nil {
		return shim.Error(err.Error())
	}

	// write the variable into the ledger
	key := args[0]
	value := args[1]
	err = stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end write")
	return shim.Success(nil)
}
