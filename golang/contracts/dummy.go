package main

import (
	"errors"
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

// P1E 1 param, return err
func (*StupidContract) P1E() error {
	return errors.New("StupidContract:PRE")
}
