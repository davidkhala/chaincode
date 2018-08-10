package golang

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"testing"
	"fmt"
	"github.com/hyperledger/fabric/common/ledger/testutil"
)

type TestChaincode struct {
}

const (
	name = "TestChaincode"
)

var logger = shim.NewLogger(name)

func (t *TestChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Init ###########")
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *TestChaincode) Invoke(ccAPI shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### " + name + " Invoke ###########")
	fcn, _ := ccAPI.GetFunctionAndParameters()
	switch fcn {
	case "panic":
		PanicString("use panic")
	}
	return shim.Success([]byte(nil))
}

var cc = new(TestChaincode)
var mock = shim.NewMockStub(name, cc)

func TestTestChaincode_Init(t *testing.T) {
	var initArgs [][]byte
	initArgs = append(initArgs, []byte("Initfcn")) //fcn
	var TxID = "ob"

	var response = mock.MockInit(TxID, initArgs)
	fmt.Println("init ", response)
	testutil.AssertSame(t, response.Status, int32(200));
}
func TestTestChaincode_Invoke(t *testing.T) {

	var args [][]byte
	args = append(args, []byte("Invokefcn")) //fcn

	var TxID = "oa"
	var response = mock.MockInvoke(TxID, args)
	fmt.Println("invoke ", response)
	testutil.AssertSame(t, response.Status, int32(200));
	//	when error status is 500
}

func TestCreateCompositeKey(t *testing.T) {
	var cKey1 = CreateCompositeKey(mock, "a", []string{"c", "C"})
	var cKey2 = CreateCompositeKey(mock, "a", []string{"d", "D"})
	var TxID = "composityKey"
	mock.MockTransactionStart(TxID);
	PutState(mock, cKey1, []byte("c"))
	PutState(mock, cKey2, []byte("C"))
	mock.MockTransactionEnd(TxID);
	TxID = "composite1"
	mock.MockTransactionStart(TxID)
	iterator := GetStateByPartialCompositeKey(mock, "a", []string{"d"})
	var kvs States
	kvs.ParseStates(iterator)
	fmt.Println(kvs)
	mock.MockTransactionEnd(TxID)
}

/**
[31m2018-07-09 12:46:27.277 HKT [mock] HasNext -> ERRO 001[0m HasNext() couldn't get Current
mockstub.go line 410:	mockLogger.Error("HasNext() couldn't get Current")
 */
func TestWorldStates(t *testing.T) {
	var TxID = "composityKey"
	mock.MockTransactionStart(TxID)

	PutState(mock, "a_1", []byte("c"))
	PutState(mock, "a_2", []byte("C"))
	PutState(mock, "a_3", []byte("C"))

	mock.MockTransactionEnd(TxID);
	TxID = "composite1"
	mock.MockTransactionStart(TxID)
	kvs := WorldStates(mock, "")

	fmt.Println(kvs)

	kvs.ParseStates(GetStateByRange(mock, "a_1", ""))
	fmt.Println(kvs)
	mock.MockTransactionEnd(TxID)
}

func TestGetStateObj(t *testing.T) {

	var value = KVJson{"1", "2", "3"}
	var TxID = "compositeKey"
	var key = "a_1"
	mock.MockTransactionStart(TxID)

	PutStateObj(mock, key, value)

	mock.MockTransactionEnd(TxID);
	TxID = "composite1"
	mock.MockTransactionStart(TxID)

	var recovered KVJson
	GetStateObj(mock, key, &recovered)

	fmt.Println(recovered)
	mock.MockTransactionEnd(TxID)
}
func TestModifyValue(t *testing.T) {
	var key = "a_1"
	var kv KVJson
	TxID := "composite2"
	mock.MockTransactionStart(TxID)
	var modifier = func(v interface{}) {
		fmt.Println("modifierTest", v)
		fmt.Println("modifierTest", kv)
	}

	ModifyValue(mock, key, modifier, &kv)
	mock.MockTransactionEnd(TxID)
}
