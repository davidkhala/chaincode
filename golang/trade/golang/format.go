package golang

import (
	"strconv"
	"encoding/json"
)

func ToInt(bytes []byte) int {
	if bytes == nil {
		return 0
	}
	i, err := strconv.Atoi(string(bytes))
	if err != nil {
		panic(err)
	}
	return i
}
func ToBytes(i int) []byte {
	return []byte(strconv.Itoa(i))
}
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

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
