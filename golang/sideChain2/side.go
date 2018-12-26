package main

import (
	. "github.com/MediConCenHK/go-chaincode-common"
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name         = "sideChaincode2"
	mainCC       = "mainChain"
	collectionTx = "tx"
)

type SideChaincode struct {
	CommonChaincode
}

func (t SideChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	t.Logger.Info(" Init ")
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t SideChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)

	t.Prepare(stub)
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Invoke", fcn, params)
	var responseBytes []byte
	switch fcn {
	case "get":
		var args = ArgsBuilder("get")
		var key = params[0]
		args.AppendArg(key)
		responseBytes = t.InvokeChaincode(mainCC, args.Get(), "").Payload
	case "putPrivate":
		var key = params[0]
		var value = params[1]
		t.PutPrivateData(collectionTx, key, []byte(value))
		var token = stub.GetTxID()
		var tokenData = TokenData{MetaData: []byte(key)}
		PutToken(t.CommonChaincode, token, tokenData)
		responseBytes = []byte(token)
	case "getPrivate":
		var token = params[0]
		var tokenData = GetToken(t.CommonChaincode, token)
		if tokenData == nil {
			PanicString("token[" + token + "] not found")
		}
		var txKey = string(tokenData.MetaData)
		responseBytes = t.GetPrivateData(collectionTx, txKey)
	}
	t.Logger.Debug("response", responseBytes)
	return shim.Success(responseBytes)
}

func main() {
	var cc = SideChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
