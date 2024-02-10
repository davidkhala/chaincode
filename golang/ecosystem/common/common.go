package common

import (
	golang "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/davidkhala/fabric-common/golang/format"
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
	MintTime     time.Time
	Issuer       format.MSPID
	IssuerClient string

	TransferTime time.Time    // latest operation time
	Owner        string       // latest owner, any format
	OwnerType    OwnerType    // latest ownerType
	Manager      format.MSPID // latest manager
	Client       string       // latest Operator
}

type TokenCreateRequest struct {
	Owner string
}
type TokenHistory struct {
	TxId string
	TokenData
	IsDelete bool
}

func (t TokenCreateRequest) Build(identity cid.ClientIdentity, c golang.CommonChaincode) TokenData {

	timestamp := c.GetTxTimestamp()
	return TokenData{
		Owner:        t.Owner,
		MintTime:     timestamp.AsTime(),
		OwnerType:    OwnerTypeMember,
		Issuer:       identity.MspID,
		Manager:      identity.MspID,
		IssuerClient: identity.GetID(),
	}
}

type TokenTransferRequest struct {
	Owner     string
	OwnerType OwnerType
}

func (data *TokenData) Apply(request TokenTransferRequest, c golang.CommonChaincode, mspid format.MSPID) *TokenData {

	data.Owner = request.Owner
	data.OwnerType = request.OwnerType
	timestamp := c.GetTxTimestamp()
	data.TransferTime = timestamp.AsTime()
	data.Manager = mspid

	return data
}

func Hash(data []byte) string {
	return goutils.HexEncode(crypto.HashSha512(data))
}
