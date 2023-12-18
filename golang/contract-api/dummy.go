package main

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

type StupidContract struct {
	contractapi.Contract
}

func (*StupidContract) Ping() {

}
func (*StupidContract) Panic() {
	panic("StupidContract")
}
