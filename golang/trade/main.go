package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/davidkhala/chaincode/golang/trade/golang"
	"strings"
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
func (cc TradeChaincode) invokeCreatorCheck(id ID) {
	if cc.Mock {
		return
	}
	var mspList golang.StringList
	var key = cc.MSPIDListKey()
	golang.GetStateObj(*cc.CCAPI, key, &mspList)
	var creator = golang.GetThisCreator(*cc.CCAPI)
	var thisMsp = creator.Msp
	var commonName = creator.Certificate.Subject.CommonName
	logger.Debug("subject common name", commonName)
	if !mspList.Has(thisMsp) {
		golang.PanicString("thisMsp:" + thisMsp + " not included in " + mspList.String())
	}
	if id.Name != commonName {
		golang.PanicString("ID.Name:" + id.Name + " mismatched with Certificate.Subject.CommonName:" + commonName)
	}
}
func (cc TradeChaincode) mspMatch(matchMSP string) {
	if cc.Mock {
		return
	}
	var thisMsp = golang.GetThisCreator(*cc.CCAPI).Msp
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
func (cc TradeChaincode) getTxKey(tt_type string) string {
	var txID = (*cc.CCAPI).GetTxID()
	var time = golang.GetTxTime(*cc.CCAPI)
	var timeMilliSecond = golang.UnixMilliSecond(time)
	return golang.ToString(timeMilliSecond) + "|" + tt_type + "|" + txID
}
func txKeyFilter(txid string, tt_type string) bool {
	var strs = strings.Split(txid, "|")
	return strs[1] == tt_type
}
func checkTo(to ID,allowedType string,transactionType string){
	if to.Type!= allowedType{
		golang.PanicString("invalid transaction target type:"+to.Type +" for transactionType:"+ transactionType)
	}
}

// Transaction makes payment of X units from A to B
func (t *TradeChaincode) Invoke(ccAPI shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Invoke ###########")
	t.Prepare(&ccAPI)
	if !t.Mock && !t.Debug {
		defer golang.PanicDefer(&response)
	}

	var fcn, params = ccAPI.GetFunctionAndParameters()
	response = shim.Success(nil)
	var txID = t.getTxKey(fcn)
	logger.Info("txID:" + txID)

	var id ID
	var inputTransaction CommonTransaction
	var filter Filter
	if len(params) == 0 {
		golang.PanicString("First arg required")
	}

	golang.FromJson([]byte(params[0]), &id)
	if len(params) > 1 {
		golang.FromJson([]byte(params[1]), &inputTransaction)
	}
	if len(params) > 2 {
		golang.FromJson([]byte(params[2]), &filter)
	}
	var filterTime = func(v interface{}) bool {
		var t = v.(golang.KeyModification).Timestamp
		return t > filter.Start && t < filter.End
	}
	var filterStatus = func(transaction PurchaseTransaction) bool {
		return filter.Status == "" || transaction.Status == filter.Status
	}
	t.invokeCreatorCheck(id)

	switch fcn {
	case fcnWalletCreate:
		var walletValue = WalletValue{txID, 0}
		var walletValueBytes []byte
		var wal = id.getWallet()
		var value = CommonTransaction{
			id, id, 0,
			fcnWalletCreate, golang.UnixMilliSecond(golang.GetTxTime(*t.CCAPI)),
		}
		if id.Type == MerchantType {
			walletValueBytes = golang.GetState(*t.CCAPI, wal.escrowID)
			if walletValueBytes != nil {
				return shim.Error("escrow Wallet " + wal.escrowID + " exist")
			}
			golang.PutStateObj(ccAPI, wal.escrowID, walletValue)
		}
		walletValueBytes = golang.GetState(*t.CCAPI, wal.regularID)
		if walletValueBytes != nil {
			return shim.Error("Wallet " + wal.regularID + " exist")
		}
		golang.PutStateObj(ccAPI, wal.regularID, walletValue)

		golang.PutStateObj(ccAPI, txID, value)
	case fcnWalletBalance:
		var walletValue WalletValue
		var wallet = t.getWalletIfExist(id)
		golang.GetStateObj(ccAPI, wallet.regularID, &walletValue)
		response = shim.Success([]byte(golang.ToString(walletValue.Balance)))
	case fcnTransfer:

		var value = CommonTransaction{
			id, inputTransaction.To, inputTransaction.Amount,
			tt, inputTransaction.TimeStamp,
		}

		var toWalletValue WalletValue
		var fromWalletValue WalletValue
		var toWallet = t.getWalletIfExist(value.To)
		var fromWallet = t.getWalletIfExist(value.From)
		golang.ModifyValue(ccAPI, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)
		golang.ModifyValue(ccAPI, fromWallet.regularID, fromWalletValue.Lose(value.Amount, txID), &fromWalletValue)
		golang.PutStateObj(ccAPI, txID, value)

	case fcnHistory:
		var wallet = t.getWalletIfExist(id)
		var historyResponse = HistoryResponse{
			wallet, nil, nil,
		}
		var filterFunc = func(v interface{}) bool {
			var t = v.(golang.KeyModification).Timestamp
			return t > filter.Start && t < filter.End
		}

		if id.Type == MerchantType {
			var escrowHistory golang.History
			var escrowHistoryIter = golang.GetHistoryForKey(ccAPI, wallet.escrowID)

			escrowHistory.ParseHistory(escrowHistoryIter, filterFunc)
			var result []CommonTransaction
			for _, entry := range escrowHistory.Modifications {
				var walletValue WalletValue
				golang.FromJson(entry.Value, &walletValue)
				var key = walletValue.RecordID
				var tx CommonTransaction
				golang.GetStateObj(ccAPI, key, &tx)
				result = append(result, tx)
			}
			historyResponse.EscrowHistory = result
		}

		var regularHistory golang.History
		var regularHistoryIter = golang.GetHistoryForKey(ccAPI, wallet.regularID)
		regularHistory.ParseHistory(regularHistoryIter, filterFunc)
		var result []CommonTransaction
		for _, entry := range regularHistory.Modifications {
			var walletValue WalletValue
			golang.FromJson(entry.Value, &walletValue)
			var key = walletValue.RecordID
			var tx CommonTransaction
			golang.GetStateObj(ccAPI, key, &tx)
			result = append(result, tx)
		}
		historyResponse.RegularHistory = result

		response = shim.Success(golang.ToJson(historyResponse))

	case tt_new_eToken_issue:
		t.mspMatch(ExchangerMSP)
		var toWallet = t.getWalletIfExist(id)

		var value = CommonTransaction{
			ID{}, id, inputTransaction.Amount,
			tt_new_eToken_issue, inputTransaction.TimeStamp,
		}

		var toWalletValue WalletValue
		golang.ModifyValue(ccAPI, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)

		golang.PutStateObj(ccAPI, txID, value)
	case tt_fiat_eToken_exchange:
		switch id.Type {
		case ConsumerType:
			t.mspMatch(ConsumerMSP)
			checkTo(inputTransaction.To,ExchangerType,tt_fiat_eToken_exchange)
		case ExchangerType:
			t.mspMatch(ExchangerMSP)
			checkTo(inputTransaction.To,ConsumerType,tt_fiat_eToken_exchange)
		default:
			golang.PanicString("invalid user type to exchange token:" + id.Type)
		}

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
		checkTo(inputTransaction.To,MerchantType,tt_consumer_purchase)
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
	case listConsumerPurchase:
		t.mspMatch(ConsumerMSP)

		var wallet = t.getWalletIfExist(id)

		var historyResponse = HistoryPurchase{
			nil,
		}

		var regularHistory golang.History
		var regularHistoryIter = golang.GetHistoryForKey(ccAPI, wallet.regularID)
		regularHistory.ParseHistory(regularHistoryIter, filterTime)
		var result []PurchaseTransaction
		for _, entry := range regularHistory.Modifications {
			var walletValue WalletValue
			golang.FromJson(entry.Value, &walletValue)
			var key = walletValue.RecordID
			if ! txKeyFilter(key, tt_consumer_purchase) {
				continue
			}
			var tx PurchaseTransaction
			golang.GetStateObj(ccAPI, key, &tx)
			if ! filterStatus(tx) {
				continue
			}
			result = append(result, tx)
		}
		historyResponse.History = result

		response = shim.Success(golang.ToJson(historyResponse))

	case listMerchantPurchase:
		t.mspMatch(MerchantMSP)

		var wallet = t.getWalletIfExist(id)

		var historyResponse = HistoryPurchase{
			nil,
		}

		var regularHistory golang.History
		var regularHistoryIter = golang.GetHistoryForKey(ccAPI, wallet.escrowID)
		regularHistory.ParseHistory(regularHistoryIter, filterTime)
		var result []PurchaseTransaction
		for _, entry := range regularHistory.Modifications {
			var walletValue WalletValue
			golang.FromJson(entry.Value, &walletValue)
			var key = walletValue.RecordID
			if ! txKeyFilter(key, tt_consumer_purchase) {
				continue
			}
			var tx PurchaseTransaction
			golang.GetStateObj(ccAPI, key, &tx)
			if ! filterStatus(tx) {
				continue
			}
			result = append(result, tx)
		}
		historyResponse.History = result

		response = shim.Success(golang.ToJson(historyResponse))

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
