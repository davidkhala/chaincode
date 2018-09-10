package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
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
	fcn, args := stub.GetFunctionAndParameters()
	logger.Info("fcn:" + fcn)
	logger.Info("args")
	for i, val := range args {
		logger.Info("arg[" + strconv.Itoa(i) + "] " + val)
	}
	return shim.Success(nil)
}

func main() {
	shim.Start(new(SideChaincode))
}
