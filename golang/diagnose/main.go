package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/fabric-common-chaincode-golang/ext"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name = "diagnose"
)

type diagnoseChaincode struct {
	CommonChaincode
}

func (t diagnoseChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Init", fcn, params)
	t.printTransient()
	return shim.Success(nil)

}
func (t diagnoseChaincode) printTransient() {
	var transientMap = t.GetTransient()
	t.Logger.Debug("==[start]transientMap")
	for k, v := range transientMap {
		t.Logger.Debug(k, ":", string(v))
	}
	t.Logger.Debug("==[end]transientMap")
}

type txData struct {
	Time  TimeLong
	Value []byte
}

func (t diagnoseChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	fcn, params := stub.GetFunctionAndParameters()
	var args = stub.GetArgs()
	t.Logger.Info("Invoke", fcn, params)
	t.printTransient()
	var responseBytes []byte
	switch fcn {
	case "panic":
		PanicString("test panic")
	case "richQuery":
		var query = params[0]
		t.Logger.Info("Query string", query)
		var queryIter = t.GetQueryResult(query)
		var states = ParseStates(queryIter, nil)
		responseBytes = ToJson(states)
	case "worldStates":
		var states = t.WorldStates("", nil)
		responseBytes = ToJson(states)
	case "whoami":
		responseBytes = ToJson(cid.NewClientIdentity(stub))
	case "get":
		var key = params[0]
		var tx txData
		t.GetStateObj(key, &tx)
		responseBytes = tx.Value
	case "putRaw":
		// for leveldb hacker analyzer
		var key = params[0]
		var value = args[2]
		t.PutState(key, value)
	case "getRaw":
		var key = params[0]
		responseBytes = t.GetState(key)
	case "history":
		var key = params[0]
		var iter = t.GetHistoryForKey(key)
		var modifications = ParseHistory(iter, nil)
		responseBytes = ToJson(modifications)
	case "put":
		var key = params[0]
		var value = params[1]
		var time TimeLong
		t.PutStateObj(key, txData{
			time.FromTimeStamp(t.GetTxTimestamp()),
			[]byte(value),
		})
	case "putEndorsement":
		var key = params[0]
		var orgs = params[1:]
		var policy = ext.NewKeyEndorsementPolicy(nil)
		policy.AddOrgs(msp.MSPRole_MEMBER, orgs...)
		t.SetStateValidationParameter(key, policy.Policy())
	case "getEndorsement":
		var key = params[0]
		var policyBytes = t.GetStateValidationParameter(key)
		var policy = ext.NewKeyEndorsementPolicy(policyBytes)
		var orgs = policy.ListOrgs()
		responseBytes = ToJson(orgs)
	case "delegate":
		type crossChaincode struct {
			ChaincodeName string
			Fcn           string
			Args          []string
			Channel       string
		}
		var paramInput crossChaincode
		FromJson([]byte(params[0]), &paramInput)
		var args = ArgsBuilder(paramInput.Fcn)
		for i, element := range paramInput.Args {
			args.AppendArg(element)
			t.Logger.Debug("delegated Arg", i, element)
		}
		var pb = t.InvokeChaincode(paramInput.ChaincodeName, args.Get(), paramInput.Channel)
		responseBytes = pb.Payload
	case "listPage":
		var startKey = params[0]
		var endKey = params[1]
		var pageSize = Atoi(params[2])
		var bookMark = params[3]
		var iter, metaData = t.GetStateByRangeWithPagination(startKey, endKey, pageSize, bookMark)

		type Response struct {
			States   []StateKV
			MetaData QueryResponseMetadata
		}
		responseBytes = ToJson(Response{ParseStates(iter, nil), metaData})
	case "GetStateByRange":
		var startKey = params[0]
		var endKey = params[1]
		var iter = t.GetStateByRange(startKey, endKey)
		responseBytes = ToJson(ParseStates(iter, nil))
	case "putBatch":
		var batch map[string]string
		FromJson([]byte(params[0]), &batch)
		for k, v := range batch {
			t.PutState(k, []byte(v))
		}
	case "chaincodeId":
		responseBytes = []byte(t.GetChaincodeID())
	case "getCertID":
		var certID = cid.NewClientIdentity(stub).GetID()
		responseBytes = []byte(certID)
	case "createComposite":
		var objectType = params[0]
		var attr1 = params[1:]
		var compositeKey = t.CreateCompositeKey(objectType, attr1)
		responseBytes = []byte(compositeKey)
	case "setEvent":
		var eventName = params[0]
		var event = params[1]
		t.SetEvent(eventName, []byte(event))
	default:
		panic("fcn not found:" + fcn)
	}
	return shim.Success(responseBytes)
}

func main() {
	var cc = diagnoseChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)
}
