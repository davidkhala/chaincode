package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type StressChaincode struct {
}

func (t *StressChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *StressChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(StressChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
