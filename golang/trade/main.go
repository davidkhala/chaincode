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
	golang.CommonChaincode
}

func (cc TradeChaincode) MSPIDListKey() string {
	return golang.CreateCompositeKey(*cc.CCAPI, MSP, []string{"ID"})
}
func (cc TradeChaincode) initMSPAllow() {
	if cc.Mock {
		return
	}
	var list = golang.StringList{
		[]string{ConsumerMSP, ExchangerMSP, MerchantMSP},
	}
	var key = cc.MSPIDListKey()
	golang.PutStateObj(*cc.CCAPI, key, list)
}
func (cc TradeChaincode) invokeMSPCheck() {
	if cc.Mock {
		return
	}
	var mspList golang.StringList
	var key = cc.MSPIDListKey()
	golang.GetStateObj(*cc.CCAPI, key, &mspList)
	var thisMsp = golang.GetThisMsp(*cc.CCAPI)
	if !mspList.Has(thisMsp) {
		golang.PanicString("thisMsp:" + thisMsp + " not included in " + mspList.String())
	}
}
func (cc TradeChaincode) mspMatch(matchMSP string) {
	if cc.Mock {
		return
	}
	var thisMsp = golang.GetThisMsp(*cc.CCAPI)
	if thisMsp != matchMSP {
		golang.PanicString("This MSP " + thisMsp + "is not allowed to operate")
	}
}

func (t *TradeChaincode) Init(ccAPI shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Init ###########")
	t.Prepare(&ccAPI)
	if !t.Mock && !t.Debug {
		defer golang.PanicDefer(&response)
	}

	t.initMSPAllow()
	response = shim.Success(nil)
	return response
}

func (cc TradeChaincode) getWalletIfExist(id ID) (wallet) {
	var walletValueBytes []byte
	var wal = id.getWallet()
	if id.Type == MerchantType {
		walletValueBytes = golang.GetState(*cc.CCAPI, wal.escrowID)
		if walletValueBytes == nil {
			golang.PanicString("escrow Wallet " + wal.escrowID + " not exist")
		}
	}
	walletValueBytes = golang.GetState(*cc.CCAPI, wal.regularID)
	if walletValueBytes == nil {
		golang.PanicString("Wallet " + wal.regularID + " not exist")
	}
	return wal
}
func (cc TradeChaincode) getPurchaseTxIfExist(purchaseTxID string) PurchaseTransaction {
	//TODO value checking with defer
	var valueBytes = golang.GetState(*cc.CCAPI, purchaseTxID)
	if valueBytes == nil {
		golang.PanicString("PurchaseTxID:" + purchaseTxID + " not exist")
	}
	var tx PurchaseTransaction
	golang.FromJson(valueBytes, &tx)
	return tx;
}
func (cc TradeChaincode) getTxKey() string {
	var txID = (*cc.CCAPI).GetTxID()
	var time = golang.GetTxTime(*cc.CCAPI)
	var timeMilliSecond = golang.UnixMilliSecond(time)
	return string(golang.ToBytes(timeMilliSecond)) + "|" + txID
}

// Transaction makes payment of X units from A to B
func (t *TradeChaincode) Invoke(ccAPI shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Invoke ###########")
	t.Prepare(&ccAPI)
	if !t.Mock && !t.Debug {
		defer golang.PanicDefer(&response)
	}
	t.invokeMSPCheck()
	var fcn, params = ccAPI.GetFunctionAndParameters()
	response = shim.Success(nil)
	var txID = t.getTxKey()
	logger.Info("txID:" + txID)

	var id ID
	if len(params) == 0 {
		golang.PanicString("First arg required")
	} else {
		golang.FromJson([]byte(params[0]), &id)
	}
	var inputTransaction CommonTransaction

	switch fcn {
	case fcnWalletCreate:
		var walletValue = WalletValue{"", 0}
		var wallet = id.getWallet()
		if id.Type == MerchantType {
			golang.PutStateObj(ccAPI, wallet.regularID, walletValue)
			golang.PutStateObj(ccAPI, wallet.escrowID, walletValue)
		} else {
			golang.PutStateObj(ccAPI, wallet.regularID, walletValue)
		}
	case fcnWalletBalance:
		var walletValue WalletValue
		var wallet = t.getWalletIfExist(id)
		golang.GetStateObj(ccAPI, wallet.regularID, &walletValue)
		response = shim.Success(golang.ToBytes(walletValue.Balance))
	case fcnHistory:
		var wallet = t.getWalletIfExist(id)
		var history = golang.ParseHistory(golang.GetHistoryForKey(ccAPI, wallet.regularID))
		var result = HistoryTransactions{[]CommonTransaction{}}
		for _, entry := range history {
			var key = entry.Value
			var tx CommonTransaction
			golang.GetStateObj(ccAPI, key, &tx)
			result.History = append(result.History, tx)
		}
		response = shim.Success(golang.ToJson(result))

	case tt_new_eToken_issue:
		t.mspMatch(ExchangerMSP)
		var toWallet = t.getWalletIfExist(id)
		golang.FromJson([]byte(params[1]), &inputTransaction)
		var value = CommonTransaction{
			ID{}, id, inputTransaction.Amount,
			tt_new_eToken_issue, inputTransaction.TimeStamp,
		}

		var toWalletValue WalletValue
		golang.ModifyValue(ccAPI, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)

		golang.PutStateObj(ccAPI, txID, value)
	case tt_fiat_eToken_exchange:
		t.mspMatch(ExchangerMSP)
		golang.FromJson([]byte(params[1]), &inputTransaction)
		var value = CommonTransaction{
			id, inputTransaction.To, inputTransaction.Amount,
			tt_fiat_eToken_exchange, inputTransaction.TimeStamp,
		}

		var toWalletValue WalletValue
		var fromWalletValue WalletValue
		var toWallet = t.getWalletIfExist(value.To)
		var fromWallet = t.getWalletIfExist(value.From)
		golang.ModifyValue(ccAPI, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)
		golang.ModifyValue(ccAPI, fromWallet.regularID, fromWalletValue.Lose(value.Amount, txID), &fromWalletValue)
		golang.PutStateObj(ccAPI, txID, value)

	case tt_consumer_purchase:
		t.mspMatch(ConsumerMSP)
		var inputTransaction PurchaseTransaction
		golang.FromJson([]byte(params[1]), &inputTransaction)
		var value = PurchaseTransaction{
			CommonTransaction{
				id, inputTransaction.To,
				inputTransaction.Amount, tt_consumer_purchase,
				inputTransaction.TimeStamp,
			},
			inputTransaction.MerchandiseCode,
			inputTransaction.MerchandiseAmount,
			inputTransaction.ConsumerDeliveryInstruction,
			StatusPending,
		}

		var toWalletValue WalletValue
		var fromWalletValue WalletValue
		var toWallet = t.getWalletIfExist(value.To)
		var fromWallet = t.getWalletIfExist(value.From)
		golang.ModifyValue(ccAPI, fromWallet.regularID, fromWalletValue.Lose(value.Amount, txID), &fromWalletValue)
		golang.ModifyValue(ccAPI, toWallet.escrowID, toWalletValue.Add(value.Amount, txID), &toWalletValue)
		golang.PutStateObj(ccAPI, txID, value)
		response = shim.Success([]byte(txID))
	case tt_merchant_accept_purchase:
		t.mspMatch(MerchantMSP)
		var inputTransaction PurchaseArbitrationTransaction
		golang.FromJson([]byte(params[1]), &inputTransaction)

		var purchaseTx = t.getPurchaseTxIfExist(inputTransaction.PurchaseTxID)
		var value = PurchaseArbitrationTransaction{
			CommonTransaction{
				id, id,
				purchaseTx.Amount, tt_merchant_accept_purchase,
				inputTransaction.TimeStamp,
			},
			true,
			inputTransaction.PurchaseTxID,
		}

		var toWalletValue WalletValue
		var fromWalletValue WalletValue
		var merchantWallet = t.getWalletIfExist(id)
		golang.ModifyValue(ccAPI, merchantWallet.escrowID, fromWalletValue.Lose(value.Amount, txID), &fromWalletValue)
		golang.ModifyValue(ccAPI, merchantWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)

		golang.ModifyValue(ccAPI, inputTransaction.PurchaseTxID, purchaseTx.Accept(), &purchaseTx)
		golang.PutStateObj(ccAPI, txID, value)

	case tt_merchant_reject_purchase:
		t.mspMatch(MerchantMSP)
		var inputTransaction PurchaseArbitrationTransaction
		golang.FromJson([]byte(params[1]), &inputTransaction)

		var purchaseTx = t.getPurchaseTxIfExist(inputTransaction.PurchaseTxID)
		var value = PurchaseArbitrationTransaction{
			CommonTransaction{
				id, purchaseTx.From,
				purchaseTx.Amount, tt_merchant_reject_purchase,
				inputTransaction.TimeStamp,
			},
			false,
			inputTransaction.PurchaseTxID,
		}

		var toWalletValue WalletValue
		var fromWalletValue WalletValue
		var fromWallet = t.getWalletIfExist(value.From)
		var toWallet = t.getWalletIfExist(value.To)
		golang.ModifyValue(ccAPI, fromWallet.escrowID, fromWalletValue.Lose(value.Amount, txID), &fromWalletValue)
		golang.ModifyValue(ccAPI, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)

		golang.ModifyValue(ccAPI, inputTransaction.PurchaseTxID, purchaseTx.Reject(), &purchaseTx)
		golang.PutStateObj(ccAPI, txID, value)
	default:
		golang.PanicString("undefined fcn:" + fcn)
	}
	return response

}

func main() {
	var cc = new(TradeChaincode)
	cc.Mock = false
	cc.Debug = true
	shim.Start(cc)
}
