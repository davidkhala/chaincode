package main

import (
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (*SmartContract) Who(context contractapi.TransactionContextInterface) string {
	return cid.NewClientIdentity(context.GetStub()).GetID()
}
