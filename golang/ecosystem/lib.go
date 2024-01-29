package main

import (
	"errors"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
)

func (t GlobalChaincode) putToken(cid cid.ClientIdentity, tokenID string, tokenData TokenData) {
	tokenData.Client = cid.GetID()
	t.common.PutStateObj(tokenID, tokenData)
}

func (t GlobalChaincode) getToken(tokenId string) *TokenData {

	var tokenData TokenData
	var exist = t.common.GetStateObj(tokenId, &tokenData)
	if exist {
		return &tokenData
	} else {
		return nil
	}
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
func panicEcosystem(Type, message string) {
	panic(errors.New("ECOSYSTEM|" + Type + "|" + message))
}
func (t GlobalChaincode) clientId() string {
	var identity = cid.NewClientIdentity(t.common.CCAPI)
	return identity.GetID()
}
