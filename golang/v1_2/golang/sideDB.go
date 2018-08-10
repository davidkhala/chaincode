package golang

import "github.com/hyperledger/fabric/core/chaincode/shim"

func (cc CommonChaincode) GetPrivateData(collection, key string) ([]byte) {
	var r, err = cc.CCAPI.GetPrivateData(collection, key)
	PanicError(err)
	return r
}
func (cc CommonChaincode) PutPrivateData(collection, key string, value []byte) {
	var err = cc.CCAPI.PutPrivateData(collection, key, value)
	PanicError(err)
}

func (cc CommonChaincode) GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (shim.StateQueryIteratorInterface) {
	var r, err = cc.CCAPI.GetPrivateDataByPartialCompositeKey(collection, objectType, keys)
	PanicError(err)
	return r
}
func (cc CommonChaincode) GetPrivateDataByRange(collection, startKey, endKey string) (shim.StateQueryIteratorInterface) {
	var r, err = cc.CCAPI.GetPrivateDataByRange(collection, startKey, endKey)

	PanicError(err)
	return r
}
func (cc CommonChaincode) GetPrivateDataQueryResult(collection, query string) (shim.StateQueryIteratorInterface) {
	var r, err = cc.CCAPI.GetPrivateDataQueryResult(collection, query);
	PanicError(err)
	return r;
}
