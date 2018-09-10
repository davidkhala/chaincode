package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name   = "mainChain"
	called = "sideChain"
)

var logger = shim.NewLogger(name)

type MainChaincode struct {
}

type Args struct {
	buff [][]byte
}

func ArgsBuilder(fcn string) (Args) {
	return Args{[][]byte{[]byte(fcn)}}
}

func (t *Args) AppendArg(str string) {
	t.buff = append(t.buff, []byte(str))
}
func (t Args) Get() [][]byte {
	return t.buff
}
func (t *MainChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)
}

func (t *MainChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Invoke ###########")
	var fcn = "invoke side"
	var args = ArgsBuilder(fcn)
	args.AppendArg("p1")

	stub.InvokeChaincode(called, args.Get(), "")

	return shim.Success(nil)
}

func main() {
	shim.Start(new(MainChaincode))
}
