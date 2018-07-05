package main

import (
	"testing"
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"github.com/hyperledger/fabric/common/ledger/testutil"
)

func TestErrorID(t *testing.T) {
	var errorID = ID{"david", "C"}
	defer testutil.AssertPanic(t, "msg")
	BuildID(golang.ToJson(errorID))
}
func TestBuildID(t *testing.T) {
	var name = "david"
	var prefix = "c"
	var errorID = BuildID(golang.ToJson(ID{name, prefix}))
	testutil.AssertEquals(t, errorID.Name, name)
	testutil.AssertEquals(t, errorID.Prefix, prefix)
}
