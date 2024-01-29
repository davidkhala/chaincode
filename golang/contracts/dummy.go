package main

import (
	"errors"
	"github.com/davidkhala/fabric-common-chaincode-golang/contract-api"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type StupidContract struct {
	contractapi.Contract
}

func (*StupidContract) Ping() {

}
func (*StupidContract) Panic() {
	panic("StupidContract")
}

// no param, return err
func (*StupidContract) Error() error {
	return errors.New("StupidContract:Error")
}
func (*StupidContract) UnUsedContext(context contractapi.TransactionContextInterface) {
	return
}

// OnlyParams is a dead function
//
//	Error: managing parameter param0. Conversion error. Value "..." was not passed in expected format []interface {}
func (*StupidContract) OnlyParams(context contractapi.TransactionContextInterface, params ...interface{}) interface{} {
	return params
}
func (*StupidContract) StringParam(context contractapi.TransactionContextInterface, p1 string) string {
	return p1
}

// StringParams is a dead function
//
//	Error: Inconsistent type in JSPB repeated field array. Got array expected object
func (*StupidContract) StringParams(context contractapi.TransactionContextInterface, p1 ...string) []string {
	return p1
}

func (*StupidContract) Defer() (err error) {
	defer contract_api.Deferred(contract_api.DefaultDeferHandler(&err))
	panic(errors.New("defer"))

}
