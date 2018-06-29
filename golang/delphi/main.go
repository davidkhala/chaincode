package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type StressChaincode struct {
}

func (t *StressChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *StressChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func main() {
	shim.Start(new(StressChaincode))
}
