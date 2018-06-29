package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"crypto/x509"
	"bytes"
	"encoding/pem"
	"errors"
	"encoding/json"
	"strings"
)

const (
	name     = "hashCC"
	superMSP = "TkMSP"
)

var logger = shim.NewLogger(name)

type HashChaincode struct {
}

type Creator struct {
	Msp            string
	CertificatePem string
	Certificate    x509.Certificate
}

type KVJson struct {
	Namespace     string
	Key           string
	CompositeKeys []string
	Value         string
}

func WorldStates(stub shim.ChaincodeStubInterface, objectType string) ([]KVJson, error) {
	var keysIterator shim.StateQueryIteratorInterface;
	var err error;
	if objectType == "" {
		keysIterator, err = stub.GetStateByRange("", "")
	} else {
		keysIterator, err = stub.GetStateByPartialCompositeKey(objectType, nil)
	}

	if err != nil {
		return nil, err
	}
	defer keysIterator.Close()

	var kvs []KVJson
	for keysIterator.HasNext() {
		kv, iterErr := keysIterator.Next()
		if iterErr != nil {
			return nil, iterErr
		}
		logger.Info(kv.Namespace, kv.Key, string(kv.Value))
		if objectType == "" {
			kvs = append(kvs, KVJson{Namespace: kv.Namespace, Key: kv.Key, Value: string(kv.Value)})
		} else {
			pKey, keys, err := stub.SplitCompositeKey(kv.Key)
			if err != nil {
				return nil, iterErr
			}
			keys = append([]string{pKey}, keys...)
			kvs = append(kvs, KVJson{Namespace: kv.Namespace, CompositeKeys: keys, Value: string(kv.Value)})
		}
	}
	return kvs, nil
}

func ParseCreator(creator []byte) (*Creator, error) {

	var msp bytes.Buffer

	var certificateBuffer bytes.Buffer
	var mspReady bool
	mspReady = false

	for i := 0; i < len(creator); i++ {
		char := creator[i]
		if char < 127 && char > 31 {
			if !mspReady {
				msp.WriteByte(char)
			} else {
				certificateBuffer.WriteByte(char)
			}
		} else if char == 10 {
			if (mspReady) {
				certificateBuffer.WriteByte(char)
			}
		} else {
			if msp.Len() > 0 {
				mspReady = true
			}

		}
	}

	block, rest := pem.Decode(certificateBuffer.Bytes())

	if rest != nil && len(rest) > 0 {
		return nil, errors.New("pem decode failed:" + string(rest))
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.New("pem decode failed:" + err.Error())
	}
	return &Creator{Msp: msp.String(), CertificatePem: certificateBuffer.String(), Certificate: *certificate}, nil

}

func (t *HashChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("###########" + name + " Init ###########")
	// GetStatus in Init will timeout request
	return shim.Success(nil)

}

func (t *HashChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	creatorBytes, err := stub.GetCreator();
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(string(creatorBytes))
	creator, err := ParseCreator(creatorBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(creator.Msp);
	logger.Info("###########" + name + " Invoke :" + "###########")
	fcn, args := stub.GetFunctionAndParameters()
	var result string
	if fcn == "worldStates" {
		kvs, err := WorldStates(stub, creator.Msp);
		if err != nil {
			return shim.Error(err.Error())
		}
		var stringArray []string;
		for _, kv := range kvs {
			jsonString, _ := json.Marshal(kv)
			jsonElement := string(jsonString)
			logger.Info(jsonElement)
			stringArray = append(stringArray, jsonElement)
		}
		jsonArray := "[" + strings.Join(stringArray, ",") + "]"
		return shim.Success([]byte(jsonArray))
	}

	compositeKey, err := stub.CreateCompositeKey(creator.Msp, []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	if fcn == "read" {
		if creator.Msp == superMSP {
			compositeKey, err = stub.CreateCompositeKey(args[1], []string{args[0]})
			if err != nil {
				return shim.Error(err.Error())
			}
		}
		resultBytes, err := stub.GetState(compositeKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		result = string(resultBytes)
	} else if fcn == "write" {
		err := stub.PutState(compositeKey, []byte(args[1]))
		if err != nil {
			return shim.Error(err.Error())
		}
		result = args[0] + "=>" + args[1]
	} else if fcn == "delete" {
		err := stub.DelState(compositeKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		result = args[0]
	} else {
		return shim.Error("invalid fcn" + fcn)
	}
	return shim.Success([]byte(result))
}

func main() {
	err := shim.Start(new(HashChaincode))
	if err != nil {
		logger.Errorf("Error starting chaincode %s: %s", name, err)
	}
}
