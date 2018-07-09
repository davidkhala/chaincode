package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/davidkhala/chaincode/golang/trade/golang"
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
func initMSPAllow(ccApi shim.ChaincodeStubInterface) {
	var _, params = ccApi.GetFunctionAndParameters()
	var p0 = []byte(params[0])
	var list golang.StringList
	golang.FromJson(p0, &list) //for checking
	var key = MSPIDListKey(ccApi)
	ccApi.PutState(key, p0)
}
func (t *TradeChaincode) invokeMSPCheck(ccApi shim.ChaincodeStubInterface) {
	var mspList golang.StringList
	var key = MSPIDListKey(ccApi)
	var statebytes = golang.GetState(ccApi, key)
	golang.FromJson(statebytes, &mspList)
	var thisMsp = golang.GetThisMsp(ccApi)
	if !mspList.Has(thisMsp) {
		golang.PanicString(thisMsp + " not included in " + mspList.String())
	}
}
func (t *TradeChaincode) Init(ccAPI shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Init ###########")
	defer golang.PanicDefer(&response)
	initMSPAllow(ccAPI)
	response = shim.Success(nil)
	return response
}

// Transaction makes payment of X units from A to B
func (t *TradeChaincode) Invoke(ccApi shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Invoke ###########")

	defer golang.PanicDefer(&response)

	var fcn, params = ccApi.GetFunctionAndParameters()

	switch fcn {
	case walletCreate:
		var id ID
		golang.FromJson([]byte(params[0]), &id)
		var suffix string
		if len(params) > 1 {
			suffix = params[1]
		}
		var wallet = id.getWallet(suffix)
		var walletValue = WalletValue{"", 0}
		golang.PutStateObj(ccApi, wallet.ID, walletValue)
	case "history":

	case tt_new_eToken_issue:
		var body CommonTransaction
		golang.FromJson([]byte(params[0]), &body)
		var txID = ccApi.GetTxID();
		var ToWallet = body.To
		logger.Info("txID:" + txID)
		golang.PutStateObj(ccApi, txID, body)

		var ToWalletBytes = golang.GetState(ccApi, ToWallet.ID)
		if ToWalletBytes == nil {
			golang.PanicString("Wallet " + ToWallet.ID + " not exist")
		}
		var toWalletValue WalletValue
		golang.FromJson(ToWalletBytes, &toWalletValue)
		toWalletValue.Balance += body.Amount
		toWalletValue.RecordID = txID
		golang.PutStateObj(ccApi, ToWallet.ID, toWalletValue)
	default:
		golang.PanicString("undefined fcn:" + fcn)
	}
	response = shim.Success(nil)
	return response

}

func main() {
	shim.Start(new(TradeChaincode))
}
