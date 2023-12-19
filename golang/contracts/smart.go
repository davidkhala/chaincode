package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (*SmartContract) Who(context contractapi.TransactionContextInterface) interface{} {
	return context.GetClientIdentity()
}
