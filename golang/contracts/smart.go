package main

import (
	"errors"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (*SmartContract) Who(context contractapi.TransactionContextInterface) string {
	return cid.NewClientIdentity(context.GetStub()).GetID()
}
func (*StupidContract) Standard(context contractapi.TransactionContextInterface, p1 string) (string, error) {
	if len(p1) == 0 {
		return "", errors.New("empty param")
	} else {
		return p1, nil
	}
}
