package main

import (
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/goutils"
	"github.com/davidkhala/goutils/crypto"
	"time"
)

type OwnerType string

const (
	OwnerTypeMember    = OwnerType("member")
	OwnerTypeClinic    = OwnerType("clinic")
	OwnerTypeNetwork   = OwnerType("network")
	OwnerTypeInsurance = OwnerType("insurance")
)

type TokenData struct {
	TokenCreateRequest
	Issuer       string // uses MSP ID in ecosystem
	Manager      string // uses MSP ID in ecosystem
	OwnerType    OwnerType
	IssuerClient string
	TransferTime time.Time
	Client       string // latest Operator Client
}

type TokenCreateRequest struct {
	Owner    string
	MintTime time.Time
}

func (t TokenCreateRequest) Build(identity cid.ClientIdentity) TokenData {
	return TokenData{
		TokenCreateRequest: t,
		OwnerType:          OwnerTypeMember,
		Issuer:             identity.MspID,
		Manager:            identity.MspID,
		IssuerClient:       identity.GetID(),
	}
}

type TokenTransferRequest struct {
	Owner        string
	OwnerType    OwnerType
	TransferTime time.Time
}

func (data *TokenData) Apply(request TokenTransferRequest) *TokenData {

	data.Owner = request.Owner
	data.OwnerType = request.OwnerType
	if data.TransferTime.Before(request.TransferTime) {
		data.TransferTime = request.TransferTime
	} else {
		panic("invalid TransferTime")
	}

	return data
}

func Hash(data []byte) string {
	return goutils.HexEncode(crypto.HashSha512(data))
}
