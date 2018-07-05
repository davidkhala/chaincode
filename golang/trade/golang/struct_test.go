package golang

import (
	"testing"
	"github.com/hyperledger/fabric/common/ledger/testutil"
	"fmt"
)

func TestSet_Has(t *testing.T) {
	var v1 = "abc"
	var v2 = "bcd"
	var v3 = "not"
	var set1 = StringList{};
	set1.Put(v1)
	set1.Put(v2);
	testutil.AssertEquals(t, set1.Has(v1), true)
	testutil.AssertEquals(t, set1.Has(v3), false)
	var list = set1.Strings;
	testutil.AssertEquals(t, len(list), 2)
}
func TestSet_JSON(t *testing.T) {
	var set1 StringList;
	var v1 = "abc"
	set1.Put(v1)
	var jsonBytes = ToJson(set1)
	fmt.Print(string(jsonBytes))

	var set2 StringList
	FromJson(jsonBytes, &set2)
	fmt.Println(set2)
}
