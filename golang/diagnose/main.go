package main

import (
	"encoding/json"
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strings"
)

const (
	name = "diagnose"
)

type diagnoseChaincode struct {
	CommonChaincode
}

func (t diagnoseChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	t.Logger.Info(" Init ")
	// GetStatus in Init will timeout request
	return shim.Success(nil)

}

func (t diagnoseChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	t.Prepare(stub)
	t.Logger.Info("Invoke :")

	fcn, args := stub.GetFunctionAndParameters()
	switch fcn {
	case "worldStates":
		var states = t.WorldStates("")
	case "getCreator":

	}

	var result string
	if fcn == "worldStates" {
		kvs, err := WorldStates(stub, creator.Msp);
		if err != nil {
			return shim.Error(err.Error())
		}
		var stringArray []string;
		for _, kv := range kvs {
			jsonString, _ := json.Marshal(kv)
			jsonElement := string(jsonString)
			logger.Info(jsonElement)
			stringArray = append(stringArray, jsonElement)
		}
		jsonArray := "[" + strings.Join(stringArray, ",") + "]"
		return shim.Success([]byte(jsonArray))
	}

	compositeKey, err := stub.CreateCompositeKey(creator.Msp, []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	if fcn == "read" {
		if creator.Msp == superMSP {
			compositeKey, err = stub.CreateCompositeKey(args[1], []string{args[0]})
			if err != nil {
				return shim.Error(err.Error())
			}
		}
		resultBytes, err := stub.GetState(compositeKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		result = string(resultBytes)
	} else if fcn == "write" {
		err := stub.PutState(compositeKey, []byte(args[1]))
		if err != nil {
			return shim.Error(err.Error())
		}
		result = args[0] + "=>" + args[1]
	} else if fcn == "delete" {
		err := stub.DelState(compositeKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		result = args[0]
	} else {
		return shim.Error("invalid fcn" + fcn)
	}
	return shim.Success([]byte(result))
}

func main() {
	var cc = diagnoseChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
