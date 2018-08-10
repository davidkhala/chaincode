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
	name = "PrivateDataCC"
)

var logger = shim.NewLogger(name)

func (t *PrivateDataCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *PrivateDataCC) Invoke(ccAPI shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Invoke ###########")


	ccAPI.GetPrivateData()
	return shim.Success(nil)
}

func main() {
	shim.Start(new(PrivateDataCC))
}
