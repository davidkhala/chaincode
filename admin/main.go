package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
	common "github.com/davidkhala/fabric-common-chaincode"
	"net/http"
)

//Error: error starting container:
// Failed to generate platform-specific docker build: Error returned from build: 1 "chaincode/input/src/github.com/admin/main.go:7:2:
// cannot find package "github.com/davidkhala/fabric-common/golang/chaincode-common" in any of:
// /opt/go/src/github.com/davidkhala/fabric-common/golang/chaincode-common (from $GOROOT)
// /chaincode/input/src/github.com/davidkhala/fabric-common/golang/chaincode-common (from $GOPATH)
// /opt/gopath/src/github.com/davidkhala/fabric-common/golang/chaincode-common
const (
	name       = "admincc"
	counterKey = "invokeCounter"
)

var logger = shim.NewLogger(name)

type AdminChaincode struct {
}

func (t *AdminChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
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

type Creator struct {
	Msp         string
	Certificate string
}

// Transaction makes payment of X units from A to B
func (t *AdminChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	stateBytes, _ := stub.GetState(counterKey)
	state := string(stateBytes)
	logger.Info("###########" + name + " Invoke :counter " + state + "###########")

	creatorBytes, _ := stub.GetCreator()

	creator, err := common.ParseCreator(creatorBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	cert := creator.Certificate
	logger.Info("issuer")
	logger.Info(cert.Issuer.CommonName)

	logger.Info("subject")
	logger.Info(cert.Subject.CommonName)

	stateInt, _ := strconv.Atoi(state)
	stateInt++
	state = strconv.Itoa(stateInt)
	stub.PutState(counterKey, []byte(state))
	return shim.Success([]byte(state))
}

func logQueryIterrator(iterator shim.StateQueryIteratorInterface) {
	for {
		if (iterator.HasNext()) {
			kv, _ := iterator.Next()
			logger.Info("kv==", kv);
		} else {
			iterator.Close();
			break;
		}
	}
}

func main() {
	err := shim.Start(new(AdminChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}