package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"errors"
)

const (
	name = "trade"
	MSP  = "MSP"
)

var logger = shim.NewLogger(name)

type TradeChaincode struct {
}

func MSPIDListKey(stub shim.ChaincodeStubInterface) string {
	return golang.CreateCompositeKey(stub, MSP, []string{"ID"})
}
func (t *TradeChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("########### " + name + " Init ###########")

	var _, params = stub.GetFunctionAndParameters();
	var p0 = []byte(params[0]);
	var list golang.StringList;
	golang.FromJson(p0, &list) //for checking
	var key = MSPIDListKey(stub)
	stub.PutState(key, p0)
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *TradeChaincode) Invoke(ccApi shim.ChaincodeStubInterface) (returned peer.Response) {
	logger.Info("########### " + name + " Invoke ###########")

	defer func() {
		if err := recover(); err != nil {
			switch x := err.(type) {
			case string:
				err = errors.New(x)
			case error:
			default:
				err = errors.New("Unknown panic")
			}
			returned.Status = shim.ERROR
			returned.Message = err.(error).Error()
		}
	}()

	var fcn, _= ccApi.GetFunctionAndParameters()

	
	var key = MSPIDListKey(ccApi)
	switch fcn {
	case "history":
		history  := golang.GetHistoryForKey(ccApi,key)
		golang.HistoryToArray(history)
	}
	var mspList golang.StringList

	var statebytes = golang.GetState(ccApi, key)
	golang.FromJson(statebytes, &mspList);
	var thisMsp = t.getThisMsp(ccApi)
	if mspList.Has(thisMsp) {
		returned = shim.Success(nil)
	} else {
		returned = shim.Error(thisMsp + " not included in " + mspList.String())
	}
	return returned

}
func (t *TradeChaincode) getThisMsp(ccApi shim.ChaincodeStubInterface) string {
	//var creatorBytes, err = ccApi.GetCreator();
	//golang.PanicError(err)
	//var creator golang.Creator
	//creator, err = golang.ParseCreator(creatorBytes)
	//golang.PanicError(err)
	//return creator.Msp;

	return "ConsumerMSP"
}

func main() {
	shim.Start(new(TradeChaincode))
}
