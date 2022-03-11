package main

import (
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/fabric-common-chaincode-golang/cid"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type GlobalChaincode struct {
	CommonChaincode
}

func (t GlobalChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)

	return shim.Success(nil)
}
func (t GlobalChaincode) putToken(cid ClientIdentity, tokenID string, tokenData TokenData) {
	tokenData.Client = cid
	t.PutStateObj(tokenID, tokenData)
}
func (t GlobalChaincode) getToken(token string) *TokenData {
	var tokenData TokenData
	var exist = t.GetStateObj(token, &tokenData)
	if !exist {
		return nil
	}
	return &tokenData
}
func (t GlobalChaincode) history(token string) []byte {
	var history = ParseHistory(t.GetHistoryForKey(token), nil)
	return ToJson(history)

}

func (t GlobalChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)

	var fcn, params = stub.GetFunctionAndParameters()
	var clientID = NewClientIdentity(stub)
	var MspID = clientID.MspID
	var responseBytes []byte
	var transient = t.GetTransient()
	var tokenRaw = transient["token"]
	if tokenRaw == nil {
		panicEcosystem("token", "param:token is empty")
	}
	var tokenID = Hash(tokenRaw)

	var tokenData TokenData
	var time TimeLong
	switch fcn {
	case FcnCreateToken:
		var createRequest TokenCreateRequest
		FromJson([]byte(params[0]), &createRequest)
		var tokenDataPtr = t.getToken(tokenID)
		if tokenDataPtr != nil {
			panicEcosystem("token", "token["+string(tokenRaw)+"] already exist")
		}
		tokenData = createRequest.Build()
		tokenData.OwnerType = OwnerTypeMember
		tokenData.TransferTime = TimeLong(0)
		tokenData.Issuer = MspID
		tokenData.Manager = MspID
		tokenData.IssuerClient = clientID
		t.putToken(clientID, tokenID, tokenData)

	case FcnGetToken:
		var tokenDataPtr = t.getToken(tokenID)
		if tokenDataPtr == nil {
			break
		}
		responseBytes = ToJson(*tokenDataPtr)
	case FcnTokenHistory:
		responseBytes = t.history(tokenID)
	case FcnDeleteToken:
		var tokenDataPtr = t.getToken(tokenID)
		if tokenDataPtr == nil {
			break // not exist, swallow
		}
		tokenData = *tokenDataPtr
		if MspID != tokenData.Manager {
			panicEcosystem("CID", "["+string(tokenRaw)+"]Token Data Manager("+tokenData.Manager+") mismatched with tx creator MspID: "+MspID)
		}
		t.DelState(tokenID)
	case FcnMoveToken:
		var transferReq TokenTransferRequest

		FromJson([]byte(params[0]), &transferReq)

		var tokenDataPtr = t.getToken(tokenID)
		if tokenDataPtr == nil {
			panicEcosystem("token", "token["+string(tokenRaw)+"] not found")
		}
		tokenData = *tokenDataPtr
		if tokenData.OwnerType != OwnerTypeMember {
			panicEcosystem("OwnerType", "original token OwnerType should be member, but got "+tokenData.OwnerType.To())
		}
		if tokenData.TransferTime != TimeLong(0) {
			panicEcosystem("token", "token["+string(tokenRaw)+"] was transferred")
		}

		tokenData.Apply(transferReq)
		tokenData.Manager = MspID
		tokenData.OwnerType = OwnerTypeNetwork
		tokenData.TransferTime = time.FromTimeStamp(t.GetTxTimestamp())
		t.putToken(clientID, tokenID, tokenData)
	default:
		panicEcosystem("unknown", "unknown fcn:"+fcn)
	}
	return shim.Success(responseBytes)
}

func main() {
	var cc = GlobalChaincode{}

	ChaincodeStart(cc)
}
