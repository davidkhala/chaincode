package golang

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func WorldStates(stub shim.ChaincodeStubInterface, objectType string) ([]KVJson, error) {
	var keysIterator shim.StateQueryIteratorInterface
	var err error;
	if objectType == "" {
		keysIterator, err = stub.GetStateByRange("", "")
	} else {
		keysIterator, err = stub.GetStateByPartialCompositeKey(objectType, nil)
	}

	if err != nil {
		return nil, err
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
			pKey, keys, err := stub.SplitCompositeKey(kv.Key)
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

func CreateCompositeKey(stub shim.ChaincodeStubInterface, first string, attrs []string) string {
	var key, err = stub.CreateCompositeKey(first, attrs)
	PanicError(err);
	return key
}

func GetState(stub shim.ChaincodeStubInterface, key string) []byte {
	var bytes, err = stub.GetState(key)
	PanicError(err);
	return bytes
}
func PutState(stub shim.ChaincodeStubInterface, key string, value []byte) {
	var err = stub.PutState(key, value)
	PanicError(err);
}
