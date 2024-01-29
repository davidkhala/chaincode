package main

import (
	golang "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/fabric-common-chaincode-golang/contract-api"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type GlobalChaincode struct {
	golang.CommonChaincode
	contractapi.Contract
}

func (t GlobalChaincode) GetToken(contractInterface contractapi.TransactionContextInterface) (data *TokenData, err error) {
	defer contract_api.Deferred(contract_api.DefaultDeferHandler(&err))
	t.Prepare(contractInterface.GetStub())

	data = t.getToken(t.tokenId())
	return
}
func (t GlobalChaincode) TokenHistory(contractInterface contractapi.TransactionContextInterface) (result string, err error) {
	defer contract_api.Deferred(contract_api.DefaultDeferHandler(&err))
	t.Prepare(contractInterface.GetStub())

	var tokenId = t.tokenId()
	var history = golang.ParseHistory(t.GetHistoryForKey(tokenId), nil)
	result = string(ToJson(history))
	return

}

func (t GlobalChaincode) CreateToken(contractInterface contractapi.TransactionContextInterface, createRequest TokenCreateRequest) (err error) {
	defer contract_api.Deferred(contract_api.DefaultDeferHandler(&err))
	t.Prepare(contractInterface.GetStub())

	var clientID = cid.NewClientIdentity(t.CCAPI)
	var tokenID = t.tokenId()
	var tokenDataPtr = t.getToken(tokenID)
	if tokenDataPtr != nil {
		panicEcosystem("token", "token["+string(t.tokenRaw())+"] already exist")
	}
	var tokenData TokenData
	tokenData = createRequest.Build(clientID, t.CommonChaincode)
	t.putToken(clientID, tokenID, tokenData)
	return
}

func (t GlobalChaincode) DeleteToken(contractInterface contractapi.TransactionContextInterface) (err error) {
	defer contract_api.Deferred(contract_api.DefaultDeferHandler(&err))
	t.Prepare(contractInterface.GetStub())

	var clientID = cid.NewClientIdentity(t.CCAPI)
	var MspID = clientID.MspID
	var tokenId = t.tokenId()
	var tokenData = t.getToken(tokenId)
	if tokenData == nil {
		return
	}
	if clientID.GetID() != tokenData.Client {
		panicEcosystem("CID", "["+string(t.tokenRaw())+"]Token Data Client("+tokenData.Client+") mismatched with tx creator ID: "+clientID.GetID())
	}
	if MspID != tokenData.Manager {
		panicEcosystem("CID", "["+string(t.tokenRaw())+"]Token Data Manager("+tokenData.Manager+") mismatched with tx creator MspID: "+MspID)
	}
	t.DelState(tokenId)
	return
}
func (t GlobalChaincode) MoveToken(contractInterface contractapi.TransactionContextInterface, transferReq TokenTransferRequest) (err error) {
	defer contract_api.Deferred(contract_api.DefaultDeferHandler(&err))
	t.Prepare(contractInterface.GetStub())

	var tokenId = t.tokenId()
	var tokenData = t.getToken(tokenId)
	var clientID = cid.NewClientIdentity(t.CCAPI)
	var MspID = clientID.MspID
	if tokenData == nil {
		panicEcosystem("token", "token["+string(t.tokenRaw())+"] not found")
	}
	tokenData.Apply(transferReq, t.CommonChaincode, MspID)
	t.putToken(clientID, tokenId, *tokenData)
	return
}

func main() {

	contract_api.Start(contract_api.NewChaincode(&GlobalChaincode{}))

}
