package main

import (
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/goutils"
	"github.com/davidkhala/goutils/crypto"
	"time"
)

type OwnerType byte

const (
	_ = iota
	OwnerTypeMember
	OwnerTypeClinic
	OwnerTypeNetwork
	OwnerTypeInsurance
)

func (t OwnerType) String() string {
	var enum = []string{"unknown", "member", "clinic", "network", "insurance"}
	return enum[t]
}

type TokenData struct {
	TokenCreateRequest
	Issuer       string // uses MSP ID in ecosystem
	Manager      string // uses MSP ID in ecosystem
	OwnerType    OwnerType
	IssuerClient cid.ClientIdentity
	TransferTime time.Time
	Client       cid.ClientIdentity // latest Operator Client
}

type TokenCreateRequest struct {
	Owner    string
	MintTime time.Time
}

func (t TokenCreateRequest) Build() TokenData {
	return TokenData{
		TokenCreateRequest: t,
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
