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

// PRE 1 param, return err
func PRE() error {
	return errors.New("StupidContract:PRE")
}
