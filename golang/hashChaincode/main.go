package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	common "github.com/davidkhala/fabric-common-chaincode/golang"
)

const (
	name = "hashCC"
)

var logger = shim.NewLogger(name)

type HashChaincode struct {
}

func (t *HashChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("###########" + name + " Init ###########")
	// GetStatus in Init will timeout request
	return shim.Success(nil)

}

func (t *HashChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// id, err := cid.New(stub)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// mspid, err := id.GetMSPID()
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// cert, err := id.GetX509Certificate()
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// logger.Info("id " + id)
	// logger.Info("mspid " + mspid)
	// logger.Info("cert " + cert)
	creatorBytes,err:= stub.GetCreator();
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(string(creatorBytes))
	creator,err:=common.ParseCreator(creatorBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(creator.Msp);
	logger.Info("###########" + name + " Invoke :" + "###########")
	fcn, args := stub.GetFunctionAndParameters()
	var result string
	if fcn == "read" {
		resultBytes, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		result = string(resultBytes[:])
	} else if fcn == "write" {
		err := stub.PutState(args[0], []byte(args[1]))
		if err != nil {
			return shim.Error(err.Error())
		}
		result = args[0] + "=>" + args[1]
	} else if fcn == "delete" {
		err := stub.DelState(args[0])
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
	err := shim.Start(new(HashChaincode))
	if err != nil {
		logger.Errorf("Error starting chaincode %s: %s", name, err)
	}
}
