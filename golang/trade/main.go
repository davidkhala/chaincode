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

func (cc TradeChaincode) OrgsMapKey() string {
	return golang.CreateCompositeKey(*cc.CCAPI, MSP, []string{"ID"})
}
func (cc TradeChaincode) invokeCreatorCheck(id ID, orgMap OrgMap) {
	if cc.Mock {
		return
	}
	var idMSP = orgMap[id.Type].MSPID
	var idOrgName = orgMap[id.Type].Name

	var creator = golang.GetThisCreator(*cc.CCAPI)
	var thisMsp = creator.Msp
	var commonName = creator.Certificate.Subject.CommonName
	logger.Debug("subject common name", commonName)
	var suffix = "@" + idOrgName
	var commonNameTrimmed = strings.TrimSuffix(commonName, suffix)
	if thisMsp != idMSP {
		golang.PanicString("thisMsp:" + thisMsp + " not Matched with MSP of ID:" + idMSP)
	}
	if id.Name != commonNameTrimmed {
		golang.PanicString("ID.Name:" + id.Name + " mismatched with Certificate.Subject.CommonName:" + commonName)
	}
}

func (t *TradeChaincode) Init(ccApi shim.ChaincodeStubInterface) (response peer.Response) {
	logger.Info("########### " + name + " Init ###########")
	t.Prepare(&ccApi)
	if !t.Mock && !t.Debug {
		defer golang.PanicDefer(&response)
	}
	var _, params = ccApi.GetFunctionAndParameters()
	if len(params) == 0 {
		golang.PanicString("OrgMap is required")
	}
	var orgMap OrgMap
	golang.FromJson([]byte(params[0]), &orgMap)

	var orgConsumer = orgMap[ConsumerType]
	if orgConsumer.Name == "" {
		golang.PanicString("Missing org consumer")
	}
	var orgExchange = orgMap[ExchangeType]
	if orgExchange.Name == "" {
		golang.PanicString("Missing org exchange")
	}
	var orgMerchant = orgMap[MerchantType]
	if orgMerchant.Name == "" {
		golang.PanicString("Missing org merchant")
	}
	orgMap = OrgMap{MerchantType: orgMerchant, ConsumerType: orgConsumer, ExchangeType: orgExchange}

	var key = t.OrgsMapKey()
	golang.PutStateObj(ccApi, key, orgMap)

	response = shim.Success(nil)
	return response
}

func (cc TradeChaincode) getWalletIfExist(id ID) (wallet) {
	var walletValue WalletValue
	var wal = id.getWallet()
	if id.Type == MerchantType {
		exist := golang.GetStateObj(*cc.CCAPI, wal.escrowID, &walletValue)
		if ! exist {
			golang.PanicString("escrow Wallet " + wal.escrowID + " not exist")
		}
	}
	exist := golang.GetStateObj(*cc.CCAPI, wal.regularID, &walletValue)
	if ! exist {
		golang.PanicString("Wallet " + wal.regularID + " not exist")
	}
	return wal
}
func (cc TradeChaincode) getPurchaseTxIfExist(purchaseTxID string) PurchaseTransaction {
	var tx PurchaseTransaction
	var exist = golang.GetStateObj(*cc.CCAPI, purchaseTxID, &tx)
	if !exist {
		golang.PanicString("PurchaseTxID:" + purchaseTxID + " not exist")
	}
	return tx;
}
func transfer(ccAPI shim.ChaincodeStubInterface, recordID string, fromID string, toID string, amount int64) {
	if fromID == toID {
		golang.PanicString("transaction operator equals to target:" + fromID)
	}
	var fromWalletValue, toWalletValue WalletValue
	golang.ModifyValue(ccAPI, fromID, fromWalletValue.Lose(amount, recordID, fromID), &fromWalletValue)
	golang.ModifyValue(ccAPI, toID, toWalletValue.Add(amount, recordID), &toWalletValue)
}
func (cc TradeChaincode) getTxKey(tt_type string) string {
	var txID = (*cc.CCAPI).GetTxID()
	var time = golang.GetTxTime(*cc.CCAPI)
	var timeMilliSecond = golang.UnixMilliSecond(time)
	return golang.ToString(timeMilliSecond) + "|" + tt_type + "|" + txID
}
func (cc TradeChaincode) checkFrom(from ID, allowedType string, transactionType string) {
	if from.Type != allowedType {
		golang.PanicString("invalid transaction operator type:" + from.Type + " for transactionType:" + transactionType)
	}
}

func checkTo(to ID, allowedType string, transactionType string) {
	if to.Type != allowedType {
		golang.PanicString("invalid transaction target type:" + to.Type + " for transactionType:" + transactionType)
	}
}
func (t *TradeChaincode) getTimeStamp() int64 {
	return golang.UnixMilliSecond(golang.GetTxTime(*t.CCAPI))
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
	var orgMap OrgMap
	var orgKey = t.OrgsMapKey()
	golang.GetStateObj(ccAPI, orgKey, &orgMap)

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
	var timeStamp = t.getTimeStamp() //inputTransaction.TimeStamp
	if len(params) > 2 {
		golang.FromJson([]byte(params[2]), &filter)
	}
	var filterTime = func(v interface{}) bool {
		var t = v.(golang.KeyModification).Timestamp
		return (filter.Start == 0 || t > filter.Start) && (t < filter.End || filter.End == 0)
	}

	var filterStatus = func(transaction PurchaseTransaction) bool {
		return filter.Status == "" || transaction.Status == filter.Status
	}
	t.invokeCreatorCheck(id, orgMap)

	switch fcn {
	case fcnWalletCreate:
		var walletValue = WalletValue{txID, 0}
		var walletValueExist WalletValue
		var wal = id.getWallet()
		var value = CommonTransaction{
			id, id, 0,
			fcnWalletCreate, timeStamp,
		}
		if id.Type == MerchantType {
			exist := golang.GetStateObj(*t.CCAPI, wal.escrowID, &walletValueExist)
			if exist {
				return shim.Error("escrow Wallet " + wal.escrowID + " exist")
			}
			golang.PutStateObj(ccAPI, wal.escrowID, walletValue)
		}
		exist := golang.GetStateObj(*t.CCAPI, wal.regularID, &walletValueExist)
		if exist {
			return shim.Error("Wallet " + wal.regularID + " exist")
		}
		golang.PutStateObj(ccAPI, wal.regularID, walletValue)

		golang.PutStateObj(ccAPI, txID, value)
	case fcnWalletBalance:
		var regularWalletValue WalletValue
		var escrowWalletValue WalletValue
		var wallet = t.getWalletIfExist(id)
		golang.GetStateObj(ccAPI, wallet.regularID, &regularWalletValue)
		if id.Type == MerchantType {
			golang.GetStateObj(ccAPI, wallet.escrowID, &escrowWalletValue)
		}

		var resp = BalanceResponse{regularWalletValue.Balance, escrowWalletValue.Balance}
		response = shim.Success(golang.ToJson(resp))
	case tt:
		var value = CommonTransaction{
			id, inputTransaction.To, inputTransaction.Amount,
			tt, timeStamp,
		}

		var toWallet = t.getWalletIfExist(value.To)
		var fromWallet = t.getWalletIfExist(value.From)
		transfer(ccAPI, txID, fromWallet.regularID, toWallet.regularID, value.Amount)
		golang.PutStateObj(ccAPI, txID, value)

	case fcnHistory:
		var wallet = t.getWalletIfExist(id)
		var historyResponse = HistoryResponse{
			id, nil, nil,
		}

		if id.Type == MerchantType {
			var escrowHistory golang.History
			var escrowHistoryIter = golang.GetHistoryForKey(ccAPI, wallet.escrowID)

			escrowHistory.ParseHistory(escrowHistoryIter, filterTime)
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
		regularHistory.ParseHistory(regularHistoryIter, filterTime)
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
		t.checkFrom(id, ExchangeType, tt_new_eToken_issue)
		var toWallet = t.getWalletIfExist(id)

		var value = CommonTransaction{
			ID{}, id, inputTransaction.Amount,
			tt_new_eToken_issue, timeStamp,
		}

		var toWalletValue WalletValue
		golang.ModifyValue(ccAPI, toWallet.regularID, toWalletValue.Add(value.Amount, txID), &toWalletValue)

		golang.PutStateObj(ccAPI, txID, value)
	case tt_fiat_eToken_exchange:
		if id.Type != ExchangeType {
			checkTo(inputTransaction.To, ExchangeType, tt_fiat_eToken_exchange)
		}

		var value = CommonTransaction{
			id, inputTransaction.To, inputTransaction.Amount,
			tt_fiat_eToken_exchange, timeStamp,
		}

		var toWallet = t.getWalletIfExist(value.To)
		var fromWallet = t.getWalletIfExist(value.From)
		transfer(ccAPI, txID, fromWallet.regularID, toWallet.regularID, value.Amount)
		golang.PutStateObj(ccAPI, txID, value)

	case tt_consumer_purchase:
		t.checkFrom(id, ConsumerType, tt_consumer_purchase)
		checkTo(inputTransaction.To, MerchantType, tt_consumer_purchase)
		var inputTransaction PurchaseTransaction
		golang.FromJson([]byte(params[1]), &inputTransaction)
		var value = PurchaseTransaction{
			CommonTransaction{
				id, inputTransaction.To,
				inputTransaction.Amount, tt_consumer_purchase,
				timeStamp,
			},
			inputTransaction.Merchandise,
			inputTransaction.ConsumerDeliveryInstruction,
			StatusPending,
		}
		value.isValid()

		var toWallet = t.getWalletIfExist(value.To)
		var fromWallet = t.getWalletIfExist(value.From)
		transfer(ccAPI, txID, fromWallet.regularID, toWallet.escrowID, value.Amount)
		golang.PutStateObj(ccAPI, txID, value)
		response = shim.Success([]byte(txID))
	case tt_merchant_accept_purchase:
		t.checkFrom(id, MerchantType, tt_merchant_accept_purchase)
		var inputTransaction PurchaseArbitrationTransaction
		golang.FromJson([]byte(params[1]), &inputTransaction)

		var purchaseTx = t.getPurchaseTxIfExist(inputTransaction.PurchaseTxID)
		var value = PurchaseArbitrationTransaction{
			CommonTransaction{
				id, id,
				purchaseTx.Amount, tt_merchant_accept_purchase,
				timeStamp,
			},
			true,
			inputTransaction.PurchaseTxID,
		}

		var merchantWallet = t.getWalletIfExist(id)
		transfer(ccAPI, txID, merchantWallet.escrowID, merchantWallet.regularID, value.Amount)

		golang.ModifyValue(ccAPI, inputTransaction.PurchaseTxID, purchaseTx.Accept(), &purchaseTx)
		golang.PutStateObj(ccAPI, txID, value)

	case tt_merchant_reject_purchase:
		t.checkFrom(id, MerchantType, tt_merchant_reject_purchase)
		var inputTransaction PurchaseArbitrationTransaction
		golang.FromJson([]byte(params[1]), &inputTransaction)

		var purchaseTx = t.getPurchaseTxIfExist(inputTransaction.PurchaseTxID)
		var value = PurchaseArbitrationTransaction{
			CommonTransaction{
				id, purchaseTx.From,
				purchaseTx.Amount, tt_merchant_reject_purchase,
				timeStamp,
			},
			false,
			inputTransaction.PurchaseTxID,
		}

		var fromWallet = t.getWalletIfExist(value.From)
		var toWallet = t.getWalletIfExist(value.To)
		transfer(ccAPI, txID, fromWallet.escrowID, toWallet.regularID, value.Amount)

		golang.ModifyValue(ccAPI, inputTransaction.PurchaseTxID, purchaseTx.Reject(), &purchaseTx)
		golang.PutStateObj(ccAPI, txID, value)
	case fcnListPurchase:
		var historyKey string
		var wallet = t.getWalletIfExist(id)
		switch id.Type {
		case ConsumerType:
			t.checkFrom(id, ConsumerType, fcnListPurchase)
			historyKey = wallet.regularID
		case MerchantType:
			t.checkFrom(id, MerchantType, fcnListPurchase)
			historyKey = wallet.escrowID
		default:
			golang.PanicString("invalid user type to view purchase list:" + id.Type)
		}

		var historyResponse = HistoryPurchase{}

		var filterTimeAndType = func(v interface{}) bool {
			var entry = v.(golang.KeyModification)
			var walletValue WalletValue
			golang.FromJson(entry.Value, &walletValue)
			var key = walletValue.RecordID
			var strs = strings.Split(key, "|")

			return filterTime(v) && strs[1] == tt_consumer_purchase
		}
		var history golang.History
		var historyIter = golang.GetHistoryForKey(ccAPI, historyKey)
		history.ParseHistory(historyIter, filterTimeAndType)
		for _, entry := range history.Modifications {
			var walletValue WalletValue
			golang.FromJson(entry.Value, &walletValue)
			var key = walletValue.RecordID
			var tx PurchaseTransaction
			golang.GetStateObj(ccAPI, key, &tx)
			if ! filterStatus(tx) {
				continue
			}

			historyResponse[key] = tx
		}

		response = shim.Success(golang.ToJson(historyResponse))

	default:
		golang.PanicString("invalid fcn:" + fcn)
	}
	return response

}

func main() {
	var cc = new(TradeChaincode)
	cc.Mock = false
	cc.Debug = false
	shim.Start(cc)
}
