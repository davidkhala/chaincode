package main

import (
	. "github.com/MediConCenHK/go-chaincode-common"
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	. "github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/goutils/crypto"

	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type tokenChaincode struct {
	CommonChaincode
}

func (t tokenChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	t.Logger.Info("Init")
	return shim.Success(nil)
}
func (t tokenChaincode) putToken(cid ClientIdentity, tokenID string, tokenData TokenData) {
	var issuer = crypto.GetDN(cid.Cert.Issuer)
	var alternativeTokenData = TokenData{
		crypto.GetDN(cid.Cert.Subject),
		issuer,
		issuer,
		0,
		0,
		0,
		cid,
		tokenData.MetaData,
	}

	t.PutStateObj(tokenID, alternativeTokenData)
}
func (t tokenChaincode) getToken(token string) *TokenData {
	var tokenData TokenData
	var exist = t.GetStateObj(token, &tokenData)
	if ! exist {
		return nil
	}
	return &tokenData
}
func (t tokenChaincode) history(tokenID string) []TokenData {
	var filter = func(modification interface{}) bool {
		return true
	}
	var history = ParseHistory(t.GetHistoryForKey(tokenID), filter)
	var result []TokenData
	for _, modification := range history {
		var tokenData TokenData
		FromJson([]byte(modification.Value), &tokenData)
		result = append(result, tokenData)
	}
	return result

}
func accessRight(identity ClientIdentity, tokenRaw string, data TokenData) {
	//TODO tune CA first
	if crypto.GetDN(identity.Cert.Issuer) != data.Manager { // allow manager to delete
		PanicString("[" + tokenRaw + "]Token Data Manager(" + data.Manager + ") mismatched with CID.Subject.CN:" + identity.Cert.Issuer.CommonName)
	}
}

func (t tokenChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)

	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Invoke:fcn", fcn)
	t.Logger.Debug("Invoke:params", params)
	var clientID = NewClientIdentity(stub)
	var responseBytes []byte
	var tokenRaw = params[0]
	if tokenRaw == "" {
		PanicString("param:token is empty")
	}
	var tokenID = Hash([]byte(tokenRaw))

	var tokenData TokenData
	switch fcn {
	case Fcn_putToken:
		FromJson([]byte(params[1]), &tokenData) //TODO test empty params
		t.putToken(clientID, tokenID, tokenData)
	case Fcn_getToken:
		tokenData = *t.getToken(tokenID)
		responseBytes = ToJson(tokenData)
	case Fcn_tokenHistory:
		responseBytes = ToJson(t.history(tokenID))
	case Fcn_deleteToken:
		var tokenDataPtr = t.getToken(tokenID)
		if tokenDataPtr == nil {
			return //not exist, swallow
		}
		tokenData = *tokenDataPtr
		accessRight(clientID, tokenRaw, tokenData)
		t.DelState(tokenID)
	case Fcn_moveToken:
		var transferReq TokenTransferRequest

		FromJson([]byte(params[1]), &transferReq)
		var tokenDataPtr = t.getToken(tokenID)
		if tokenDataPtr == nil {
			PanicString("token not found:" + tokenRaw)
		}
		tokenData = *tokenDataPtr
		accessRight(clientID, tokenRaw, tokenData)
		tokenData = transferReq.ApplyOn(tokenData)
		t.putToken(clientID, tokenID, tokenData)
	default:
		PanicString("unknown fcn:" + fcn)
	}
	t.Logger.Debug("response", string(responseBytes))
	return shim.Success(responseBytes)
}

func main() {
	var cc = tokenChaincode{}
	cc.SetLogger(GlobalCCID)
	shim.Start(cc)
}
