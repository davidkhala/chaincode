package golang

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"fmt"
)

func WorldStates(ccAPI shim.ChaincodeStubInterface, objectType string) ([]KVJson, error) {
	var keysIterator shim.StateQueryIteratorInterface
	if objectType == "" {
		keysIterator = GetStateByRange(ccAPI, "", "")
	} else {
		keysIterator = GetStateByPartialCompositeKey(ccAPI, objectType, nil)
	}

	defer keysIterator.Close()

	var kvs []KVJson
	for keysIterator.HasNext() {
		kv, iterErr := keysIterator.Next()
		if iterErr != nil {
			return nil, iterErr
		}
		if objectType == "" {
			kvs = append(kvs, KVJson{Namespace: kv.Namespace, Key: kv.Key, Value: string(kv.Value)})
		} else {
			pKey, keys, err := ccAPI.SplitCompositeKey(kv.Key)
			if err != nil {
				return nil, iterErr
			}
			keys = append([]string{pKey}, keys...)
			kvs = append(kvs, KVJson{Namespace: kv.Namespace, CompositeKeys: keys, Value: string(kv.Value)})
		}
	}
	return kvs, nil
}

type KVJson struct {
	Namespace     string
	Key           string
	CompositeKeys []string
	Value         string
}

func CreateCompositeKey(ccAPI shim.ChaincodeStubInterface, first string, attrs []string) string {
	var key, err = ccAPI.CreateCompositeKey(first, attrs)
	PanicError(err);
	return key
}

func GetState(ccAPI shim.ChaincodeStubInterface, key string) []byte {
	var bytes, err = ccAPI.GetState(key)
	PanicError(err);
	return bytes
}
func PutState(ccAPI shim.ChaincodeStubInterface, key string, value []byte) {
	var err = ccAPI.PutState(key, value)
	PanicError(err);
}

func HistoryToArray(iterator shim.HistoryQueryIteratorInterface) (result []queryresult.KeyModification) {
	defer iterator.Close()
	for {
		if iterator.HasNext() {
			keyModification, err := iterator.Next()
			PanicError(err);
			result = append(result, *keyModification)
			fmt.Print(*keyModification)
		}else {
			break
		}
	}
	return result;
}
func GetHistoryForKey(ccAPI shim.ChaincodeStubInterface, key string) (shim.HistoryQueryIteratorInterface) {
	var r, err = ccAPI.GetHistoryForKey(key)
	PanicError(err)
	return r;
}
func GetStateByPartialCompositeKey(ccAPI shim.ChaincodeStubInterface, objectType string, keys []string) shim.StateQueryIteratorInterface {
	var r, err = ccAPI.GetStateByPartialCompositeKey(objectType, keys)
	PanicError(err)
	return r
}
func GetStateByRange(ccAPI shim.ChaincodeStubInterface, startKey string, endKey string) shim.StateQueryIteratorInterface {
	var r, err = ccAPI.GetStateByRange(startKey, endKey)
	PanicError(err)
	return r
}
