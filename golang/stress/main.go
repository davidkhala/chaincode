package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "stress"
)

var logger = shim.NewLogger(name)

type StressChaincode struct {
}

func (t *StressChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *StressChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### " + name + " Invoke ###########")
	return shim.Success(nil)
}

func main() {
	shim.Start(new(StressChaincode))
}
