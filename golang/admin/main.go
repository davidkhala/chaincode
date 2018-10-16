package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "admincc"
)

type AdminChaincode struct {
	CommonChaincode
}

type couchDBValue struct {
	Time TimeLong
}

func (t AdminChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	t.Logger.Info(" Init ")
	// GetStatus in Init will timeout request

	var txId = t.CCAPI.GetTxID()
	var txTime = UnixMilliSecond(t.GetTxTime())
	t.PutStateObj(txId, couchDBValue{txTime})

	response = shim.Success(nil)
	return
}

// Transaction makes payment of X units from A to B
func (t AdminChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info(" Invoke fcn:" + fcn)
	switch fcn {
	case "panic":
		PanicString("test panic")
	case "richQuery":
		var query = params[0]
		t.Logger.Debug("Query string:" + query)
		var queryIter = t.GetQueryResult(query)
		var states = ParseStates(queryIter)
		for i, v := range states.States {
			t.Logger.Info(i, ":", v)
		}
	case "set":
		var txId = t.CCAPI.GetTxID()
		var txTime = UnixMilliSecond(t.GetTxTime())
		t.PutStateObj(txId, couchDBValue{txTime})
	}
	{
		transientMap, _ := stub.GetTransient()
		t.Logger.Info("==transientMap")
		for k, v := range transientMap {
			t.Logger.Info(k, ":", string(v))
		}
	}
	var txID = t.CCAPI.GetTxID()
	var txTime = UnixMilliSecond(t.GetTxTime())

	t.PutStateObj(txID, couchDBValue{txTime})
	return shim.Success(nil)
}

func main() {
	var cc = AdminChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
