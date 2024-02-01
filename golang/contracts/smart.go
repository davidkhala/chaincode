package main

import (
	"errors"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
)

type SmartContract struct {
	contractapi.Contract
}

func (*SmartContract) GetEvaluateTransactions() []string {
	return []string{"Who", "Standard", "Now"} // Must be Heading with UpperCase
}

func (*SmartContract) Who(context contractapi.TransactionContextInterface) string {
	return cid.NewClientIdentity(context.GetStub()).GetID()
}
func (*SmartContract) Standard(context contractapi.TransactionContextInterface, p1 string) (string, error) {
	if len(p1) == 0 {
		return "", errors.New("empty param")
	} else {
		return p1, nil
	}
}

// Now It can be serialized in node client in format of `2024-01-29T01:55:45.000Z`
func (*SmartContract) Now() time.Time {
	return time.Now()
}
