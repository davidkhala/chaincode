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
	mock bool
}

func MSPIDListKey(stub shim.ChaincodeStubInterface) string {
	return golang.CreateCompositeKey(stub, MSP, []string{"ID"})
}
func (cc TradeChaincode) initMSPAllow(ccApi shim.ChaincodeStubInterface) {
	if cc.mock {
		return
	}
	var _, params = ccApi.GetFunctionAndParameters()
	var p0 = []byte(params[0])
	var list golang.StringList
	golang.FromJson(p0, &list) //for checking
	var key = MSPIDListKey(ccApi)
	ccApi.PutState(key, p0)
}
func (cc TradeChaincode) invokeMSPCheck(ccApi shim.ChaincodeStubInterface) {
	if cc.mock {
		return
	}
	var mspList golang.StringList
	var key = MSPIDListKey(ccApi)
	var statebytes = golang.GetState(ccApi, key)
	golang.FromJson(statebytes, &mspList)
	var thisMsp = golang.GetThisMsp(ccApi)
	if !mspList.Has(thisMsp) {
		golang.PanicString("thisMsp:" + thisMsp + " not included in " + mspList.String())
	}
}
func (cc TradeChaincode) mspMatch(ccApi shim.ChaincodeStubInterface, matchMSP string) {
	if cc.mock {
		return
	}
	var thisMsp = golang.GetThisMsp(ccApi)
	if thisMsp != matchMSP {
		golang.PanicString("This MSP " + thisMsp + "is not allowed to operate")
	}
}

func (t *TradeChaincode) Init(ccAPI shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Init ###########")
	if !t.mock {
		defer golang.PanicDefer(&response)
	}

	t.initMSPAllow(ccAPI)
	response = shim.Success(nil)
	return response
}

func getWalletIfExist(ccApi shim.ChaincodeStubInterface, id ID) (wallet) {
	var walletValueBytes []byte
	var wal = id.getWallet()
	if id.Type == MerchantType {
		walletValueBytes = golang.GetState(ccApi, wal.escrowID)
		if walletValueBytes == nil {
			golang.PanicString("escrow Wallet " + wal.escrowID + " not exist")
		}
	}
	walletValueBytes = golang.GetState(ccApi, wal.regularID)
	if walletValueBytes == nil {
		golang.PanicString("Wallet " + wal.regularID + " not exist")
	}
	return wal
}
func getTxKey(ccApi shim.ChaincodeStubInterface) string {
	var txID = ccApi.GetTxID()
	var time = golang.GetTxTime(ccApi)
	var timeMilliSecond = golang.UnixMilliSecond(time)
	return string(golang.ToBytes(timeMilliSecond)) + "|" + txID
}

// Transaction makes payment of X units from A to B
func (t *TradeChaincode) Invoke(ccApi shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Invoke ###########")

	if !t.mock {
		defer golang.PanicDefer(&response)
	}
	t.invokeMSPCheck(ccApi)
	var fcn, params = ccApi.GetFunctionAndParameters()
	response = shim.Success(nil)
	var txID = getTxKey(ccApi)
	logger.Info("txID:" + txID)

	var id ID
	if len(params) == 0 {
		golang.PanicString("First arg required")
	} else {
		golang.FromJson([]byte(params[0]), &id)
	}
	var input CommonTransaction

	switch fcn {
	case walletCreate:
		var walletValue = WalletValue{"", 0}
		var wallet = id.getWallet()
		if id.Type == MerchantType {
			golang.PutStateObj(ccApi, wallet.regularID, walletValue)
			golang.PutStateObj(ccApi, wallet.escrowID, walletValue)
		} else {
			golang.PutStateObj(ccApi, wallet.regularID, walletValue)
		}
	case walletBalance:
		var walletValue WalletValue
		var wallet = getWalletIfExist(ccApi, id)
		golang.GetStateObj(ccApi, wallet.regularID, &walletValue)
		response = shim.Success(golang.ToBytes(walletValue.Balance))
	case "history":
		var wallet = getWalletIfExist(ccApi, id)
		var history = golang.ParseHistory(golang.GetHistoryForKey(ccApi, wallet.regularID))
		var result = HistoryTransactions{[]CommonTransaction{}}
		for _, entry := range history {
			var key = entry.Value
			var tx CommonTransaction
			golang.GetStateObj(ccApi, key, &tx)
			result.History = append(result.History, tx)
		}
		response = shim.Success(golang.ToJson(result))

	case tt_new_eToken_issue:
		t.mspMatch(ccApi, ExchangerMSP)
		var toWallet = getWalletIfExist(ccApi, id)
		golang.FromJson([]byte(params[1]), &input)
		var value = CommonTransaction{
			ID{}, id, input.Amount,
			tt_new_eToken_issue, input.TimeStamp,
		}

		var toWalletValue WalletValue
		golang.ModifyValue(ccApi, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)

		golang.PutStateObj(ccApi, txID, value)
	case tt_fiat_eToken_exchange:
		t.mspMatch(ccApi, ExchangerMSP)
		golang.FromJson([]byte(params[1]), &input)
		var value = CommonTransaction{
			id, input.To, input.Amount,
			tt_fiat_eToken_exchange, input.TimeStamp,
		}

		var toWalletValue WalletValue
		var fromWalletValue WalletValue
		var toWallet = getWalletIfExist(ccApi, value.To)
		var fromWallet = getWalletIfExist(ccApi, value.From)
		golang.ModifyValue(ccApi, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)
		golang.ModifyValue(ccApi, fromWallet.regularID, fromWalletValue.Lose(value.Amount, txID), &fromWalletValue)
		golang.PutStateObj(ccApi, txID, value)

	default:
		golang.PanicString("undefined fcn:" + fcn)
	}
	return response

}

func main() {
	var cc = new(TradeChaincode)
	cc.mock = false
	shim.Start(cc)
}
