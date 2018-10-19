package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type PrivateDataCC struct {
	CommonChaincode
}

const (
	name       = "PrivateDataCC"
	collection = "private1"
)

var logger = shim.NewLogger(name)

func (t *PrivateDataCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *PrivateDataCC) Invoke(ccAPI shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Invoke ###########")
	//defer golang.PanicDefer(&response)
	t.Prepare(ccAPI)
	var fcn, params = ccAPI.GetFunctionAndParameters()
	var responseBytes []byte
	switch fcn {
	case "put":
		var CN = t.GetThisCreator().Certificate.Subject.CommonName
		var txTime = UnixMilliSecond(t.GetTxTime()).String()
		t.PutPrivateData(collection, collection, []byte(CN+"|"+txTime))
	case "get":
		var pData = t.GetPrivateData(collection, collection)
		logger.Info("pData" + string(pData))
	case "put2":
		var txTime = UnixMilliSecond(t.GetTxTime()).String()
		var key2 = txTime + " 1"
		t.PutState(txTime, []byte(txTime))
		t.PutState(key2, []byte(key2))
		responseBytes = []byte(key2)
	case "get2":
		var key = params[0]
		responseBytes = t.GetState(key)
	default:

	}

	return shim.Success(responseBytes)
}

func main() {
	shim.Start(new(PrivateDataCC))
}
