package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "sideChaincode"
)

var logger = shim.NewLogger(name)

type SideChaincode struct {
}

func (t *SideChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SideChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Invoke ###########")
	response, err := stub.GetArgsSlice()
	if err != nil {
		panic(err)
	}
	return shim.Success(response)
}

func main() {
	shim.Start(new(SideChaincode))
}
