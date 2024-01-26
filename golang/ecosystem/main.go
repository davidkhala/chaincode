package main

import (
	golang "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/fabric-common-chaincode-golang/contract-api"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type GlobalChaincode struct {
	common golang.CommonChaincode
	contractapi.Contract
}

func (t GlobalChaincode) putToken(cid cid.ClientIdentity, tokenID string, tokenData TokenData) {
	tokenData.Client = cid
	t.common.PutStateObj(tokenID, tokenData)
}

func (t GlobalChaincode) GetToken(contractInterface contractapi.TransactionContextInterface) *TokenData {
	t.common.Prepare(contractInterface.GetStub())
	var tokenID = t.tokenId()
	var tokenData TokenData
	var exist = t.common.GetStateObj(tokenID, &tokenData)
	if !exist {
		return nil
	}
	return &tokenData
}
func (t GlobalChaincode) TokenHistory(contractInterface contractapi.TransactionContextInterface) string {
	t.common.Prepare(contractInterface.GetStub())
	var tokenId = t.tokenId()
	var history = golang.ParseHistory(t.common.GetHistoryForKey(tokenId), nil)
	return string(ToJson(history))

}

func (t GlobalChaincode) CreateToken(contractInterface contractapi.TransactionContextInterface, createRequest TokenCreateRequest) {
	t.common.Prepare(contractInterface.GetStub())
	var clientID = cid.NewClientIdentity(t.common.CCAPI)
	var MspID = clientID.MspID
	var tokenID = t.tokenId()
	var tokenDataPtr = t.GetToken(contractInterface)
	if tokenDataPtr != nil {
		panicEcosystem("token", "token["+string(t.tokenRaw())+"] already exist")
	}
	var tokenData TokenData
	tokenData = createRequest.Build()
	tokenData.OwnerType = OwnerTypeMember
	tokenData.Issuer = MspID
	tokenData.Manager = MspID
	tokenData.IssuerClient = clientID
	t.putToken(clientID, tokenID, tokenData)
}

func (t GlobalChaincode) tokenRaw() []byte {
	var transient = t.common.GetTransient()
	return transient["token"]
}
func (t GlobalChaincode) tokenId() string {
	var tokenRaw = t.tokenRaw()
	if tokenRaw == nil {
		panicEcosystem("token", "param:token is empty")
	}
	return Hash(tokenRaw)
}
func (t GlobalChaincode) clientId() string {
	var identity = cid.NewClientIdentity(t.common.CCAPI)
	return identity.GetID()
}

func (t GlobalChaincode) DeleteToken(contractInterface contractapi.TransactionContextInterface) {
	t.common.Prepare(contractInterface.GetStub())
	var clientID = cid.NewClientIdentity(t.common.CCAPI)
	var MspID = clientID.MspID
	var tokenData = t.GetToken(contractInterface)
	if tokenData == nil {
		return // not exist, skip
	}
	if MspID != tokenData.Manager {
		panicEcosystem("CID", "["+string(t.tokenRaw())+"]Token Data Manager("+tokenData.Manager+") mismatched with tx creator MspID: "+MspID)
	}
	t.common.DelState(t.tokenId())
}
func (t GlobalChaincode) MoveToken(contractInterface contractapi.TransactionContextInterface, transferReq TokenTransferRequest) {
	var tokenData = t.GetToken(contractInterface)
	var clientID = cid.NewClientIdentity(t.common.CCAPI)
	var MspID = clientID.MspID
	if tokenData == nil {
		panicEcosystem("token", "token["+string(t.tokenRaw())+"] not found")
	}
	if tokenData.OwnerType != OwnerTypeMember {
		panicEcosystem("OwnerType", "original token OwnerType should be member, but got "+tokenData.OwnerType.String())
	}
	if !tokenData.TransferTime.IsZero() {
		panicEcosystem("token", "token["+string(t.tokenRaw())+"] was transferred")
	}

	tokenData.Apply(transferReq)
	tokenData.Manager = MspID
	tokenData.OwnerType = OwnerTypeNetwork
	var time = t.common.GetTxTimestamp()
	tokenData.TransferTime = time.AsTime()
	t.putToken(clientID, t.tokenId(), *tokenData)
}

func main() {

	var chaincode = contract_api.NewChaincode(&GlobalChaincode{})
	contract_api.Start(chaincode)

}

func panicEcosystem(Type, message string) {
	PanicString("ECOSYSTEM|" + Type + "|" + message)
}
