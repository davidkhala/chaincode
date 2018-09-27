package main

import (
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name       = "admincc"
	counterKey = "invokeCounter"
)

var logger = shim.NewLogger(name)

type AdminChaincode struct {
}

func (t *AdminChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("###########" + name + " Init ###########")
	// GetStatus in Init will timeout request
	err := stub.PutState(counterKey, []byte("0"))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *AdminChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	stateBytes, _ := stub.GetState(counterKey)
	state := string(stateBytes)
	logger.Info("###########" + name + " Invoke :counter " + state + "###########")

	trasientMap, _ := stub.GetTransient()
	logger.Info("==transientMap")
	for k, v := range trasientMap {
		logger.Info(k, ":", string(v))
	}
	stateInt, _ := strconv.Atoi(state)
	stateInt++
	state = strconv.Itoa(stateInt)
	stub.PutState(counterKey, []byte(state))
	return shim.Success([]byte(state))
}

func main() {
	err := shim.Start(new(AdminChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
