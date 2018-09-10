package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "mainChaincode"
	called = "sideChaincode"
)

var logger = shim.NewLogger(name)

type MainChaincode struct {
}

type Args [][]byte
func ArgsBuilder(fcn string)(Args){
	return Args{[]byte(fcn)}
}
func (t Args) Get() [][]byte {
	return [][]byte(t)
}
func (t *Args) AppendArg(str string){
	var buff = *t
	buff = append(buff, []byte(str))
	t = &buff
}
func (t *MainChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)
}

func (t *MainChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Invoke ###########")
	var fcn ="invoke side"
	var args = ArgsBuilder(fcn)
	args.AppendArg("p1")

	var sideReponse = stub.InvokeChaincode(called,args.Get(),"")
	logger.Info(sideReponse)
	return shim.Success(nil)
}

func main() {
	shim.Start(new(MainChaincode))
}
