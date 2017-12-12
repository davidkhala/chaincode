package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "delphi_cc"
)

var logger = shim.NewLogger(name)

type DelphiChaincode struct {
}

func (t *DelphiChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *DelphiChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### " + name + " Invoke ###########")

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(DelphiChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
