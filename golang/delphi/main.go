package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type DelphiChaincode struct {
}

const (
	name       = "DelphiChaincode"
	counterKey = "invokeCounter"
)

var logger = shim.NewLogger(name)

func (t *DelphiChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Init ###########")
	stub.PutState(counterKey, []byte("0"))
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *DelphiChaincode) Invoke(ccAPI shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Invoke ###########")

	stateBytes, _ := ccAPI.GetState(counterKey)
	state := string(stateBytes)

	stateInt, _ := strconv.Atoi(state)
	stateInt++
	state = strconv.Itoa(stateInt)
	ccAPI.PutState(counterKey, []byte(state))
	logger.Info("###########" + name + " Invoke :counter " + state + "###########")
	return shim.Success([]byte(state))
}

func main() {
	shim.Start(new(DelphiChaincode))
}
