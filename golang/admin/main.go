package main

import (
	"net/http"
	"strconv"

	common "github.com/davidkhala/fabric-common-chaincode/golang"
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
	resp, err := http.Get("http://www.google.com/")
	if err != nil {
		logger.Info(err.Error())
	}
	if resp != nil {
		logger.Info(resp.Status)
	}
	return shim.Success(nil)

}

// Transaction makes payment of X units from A to B
func (t *AdminChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	stateBytes, _ := stub.GetState(counterKey)
	state := string(stateBytes)
	logger.Info("###########" + name + " Invoke :counter " + state + "###########")

	creatorBytes, _ := stub.GetCreator()

	creator, err := common.ParseCreator(creatorBytes)
	if err != nil {
		logger.Error(err.Error())
		logger.Error(string(creatorBytes[:]))
	} else {
		cert := creator.Certificate
		logger.Info("issuer")
		logger.Info(cert.Issuer.CommonName)

		logger.Info("subject")
		logger.Info(cert.Subject.CommonName)
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
