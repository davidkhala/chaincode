package main

import (
	"errors"
	golang "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/fabric-common-chaincode-golang/contract-api"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type GlobalChaincode struct {
	contractapi.Contract
}

func (t GlobalChaincode) putToken(cid cid.ClientIdentity, tokenID string, tokenData TokenData) {
	tokenData.Client = cid.GetID()
	t.common.PutStateObj(tokenID, tokenData)
}

func (t GlobalChaincode) getToken(tokenId string) (*TokenData, error) {

	var tokenData TokenData
	var exist = t.common.GetStateObj(tokenId, &tokenData)
	if !exist {
		return nil, nil
	} else {
		return &tokenData, nil
	}
}
func (t GlobalChaincode) GetToken(contractInterface contractapi.TransactionContextInterface) (*TokenData, error) {
	t.common.Prepare(contractInterface.GetStub())
	tokenId, err := t.tokenId()
	if err != nil {
		return nil, err
	}
	return t.getToken(tokenId)
}
func (t GlobalChaincode) TokenHistory(contractInterface contractapi.TransactionContextInterface) (string, error) {
	t.common.Prepare(contractInterface.GetStub())
	tokenId, err := t.tokenId()
	if err != nil {
		return "", err
	}
	var history = golang.ParseHistory(t.common.GetHistoryForKey(tokenId), nil)
	return string(ToJson(history)), nil

}

func (t GlobalChaincode) CreateToken(contractInterface contractapi.TransactionContextInterface, createRequest TokenCreateRequest) error {
	t.common.Prepare(contractInterface.GetStub())

	var clientID = cid.NewClientIdentity(t.common.CCAPI)
	tokenID, err := t.tokenId()
	if err != nil {
		return err
	}
	tokenDataPtr, err := t.getToken(tokenID)
	if err != nil {
		return err
	}
	if tokenDataPtr != nil {
		return panicEcosystem("token", "token["+string(t.tokenRaw())+"] already exist")
	}
	var tokenData TokenData
	tokenData = createRequest.Build(clientID)
	t.putToken(clientID, tokenID, tokenData)
	return nil
}

func (t GlobalChaincode) tokenRaw() []byte {
	var transient = t.common.GetTransient()
	return transient["token"]
}
func (t GlobalChaincode) tokenId() (string, error) {
	var tokenRaw = t.tokenRaw()
	if tokenRaw == nil {
		return "", panicEcosystem("token", "param:token is empty")
	}
	return Hash(tokenRaw), nil
}
func (t GlobalChaincode) clientId() string {
	var identity = cid.NewClientIdentity(t.common.CCAPI)
	return identity.GetID()
}

func (t GlobalChaincode) DeleteToken(contractInterface contractapi.TransactionContextInterface) error {
	t.common.Prepare(contractInterface.GetStub())
	var clientID = cid.NewClientIdentity(t.common.CCAPI)
	var MspID = clientID.MspID
	tokenId, err := t.tokenId()
	if err != nil {
		return err
	}
	tokenData, err := t.getToken(tokenId)
	if err != nil {
		return err
	}
	if tokenData == nil {
		return nil
	}
	if MspID != tokenData.Manager {
		return panicEcosystem("CID", "["+string(t.tokenRaw())+"]Token Data Manager("+tokenData.Manager+") mismatched with tx creator MspID: "+MspID)
	}

	t.common.DelState(tokenId)
	return nil
}
func (t GlobalChaincode) MoveToken(contractInterface contractapi.TransactionContextInterface, transferReq TokenTransferRequest) error {
	tokenData, err := t.GetToken(contractInterface)
	if err != nil {
		return err
	}
	var clientID = cid.NewClientIdentity(t.common.CCAPI)
	var MspID = clientID.MspID
	if tokenData == nil {
		return panicEcosystem("token", "token["+string(t.tokenRaw())+"] not found")
	}
	if tokenData.OwnerType != OwnerTypeMember {
		return panicEcosystem("OwnerType", "original token OwnerType should be member, but got "+string(tokenData.OwnerType))
	}
	if !tokenData.TransferTime.IsZero() {
		return panicEcosystem("token", "token["+string(t.tokenRaw())+"] was transferred")
	}

	tokenData.Apply(transferReq)
	tokenData.Manager = MspID
	tokenData.OwnerType = OwnerTypeNetwork
	var time = t.common.GetTxTimestamp()
	tokenData.TransferTime = time.AsTime()
	tokenId, err := t.tokenId()
	if err != nil {
		return err
	}
	t.putToken(clientID, tokenId, *tokenData)
	return nil
}

func main() {

	var chaincode = contract_api.NewChaincode(&GlobalChaincode{})
	contract_api.Start(chaincode)

}

func panicEcosystem(Type, message string) error {
	return errors.New("ECOSYSTEM|" + Type + "|" + message)
}
