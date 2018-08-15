package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/davidkhala/chaincode/golang/v1_2/golang"
)

type PrivateDataCC struct {
	golang.CommonChaincode
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
	var fcn, _ = ccAPI.GetFunctionAndParameters()
	switch fcn {
	case "put":
		var CN = t.GetThisCreator().Certificate.Subject.CommonName
		var txTime = golang.UnixMilliSecond(t.GetTxTime()).ToString()
		t.PutPrivateData(collection, collection, []byte(CN+"|"+txTime))
	case "get":
		var pData = t.GetPrivateData(collection, collection)
		logger.Info("pData" + string(pData))
	default:

	}

	return shim.Success(nil)
}

func main() {
	shim.Start(new(PrivateDataCC))
}
