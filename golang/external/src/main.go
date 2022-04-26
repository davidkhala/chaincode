package main

import (
	"github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"os"
)

type SmartContract struct {
	contractapi.Contract
}

type QueryResult struct {
	Key   string
	Value []byte
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	// range query with empty string for startKey and endKey does an open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		queryResult := QueryResult{Key: queryResponse.Key, Value: queryResponse.Value}
		results = append(results, queryResult)
	}

	return results, nil
}
func main() {

	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	goutils.PanicError(err)

	server := &shim.ChaincodeServer{
		CCID:    os.Getenv("CORE_CHAINCODE_ID_NAME"),
		Address: "0.0.0.0:9999", //os.Getenv("CHAINCODE_SERVER_ADDRESS"),
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	err = server.Start()
	goutils.PanicError(err)
}
