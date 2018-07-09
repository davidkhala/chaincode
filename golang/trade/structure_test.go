package main

import (
	"testing"
	"github.com/hyperledger/fabric/common/ledger/testutil"
)

func TestErrorID(t *testing.T) {
	var errorID = ID{"david", "C"}
	defer testutil.AssertPanic(t, "msg")
	errorID.getWallet("")
}
