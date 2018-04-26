package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"crypto/x509"
	"bytes"
	"encoding/pem"
	"errors"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	"encoding/json"
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

func WorldStates(stub shim.ChaincodeStubInterface) ([]queryresult.KV, error) {
	keysIterator, err := stub.GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer keysIterator.Close()

	var kvs []queryresult.KV
	for keysIterator.HasNext() {
		kv, iterErr := keysIterator.Next()
		if iterErr != nil {
			return nil, iterErr
		}
		kvs = append(kvs, *kv)
	}
	return kvs, nil
}
func KV2Json(kv queryresult.KV) (string) {
	type KVJson struct {
		Namespace string
		Key       string
		Value     string
	}
	kvJson:= KVJson{Namespace:kv.Namespace,Key:kv.Key,Value:string(kv.Value)}
	jsonString,_:=  json.Marshal(kvJson)
	return string(jsonString)
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

func (t *HashChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("###########" + name + " Init ###########")
	// GetStatus in Init will timeout request
	return shim.Success(nil)

}

func (t *HashChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// id, err := cid.New(stub)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// mspid, err := id.GetMSPID()
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// cert, err := id.GetX509Certificate()
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// logger.Info("id " + id)
	// logger.Info("mspid " + mspid)
	// logger.Info("cert " + cert)
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
	accessor := []string{creator.Msp};
	compositKey, err := stub.CreateCompositeKey(args[0], accessor)
	if err != nil {
		return shim.Error(err.Error())
	}
	if fcn == "read" {
		if creator.Msp == superMSP {
			compositKey, err = stub.CreateCompositeKey(args[0], []string{args[1]})
			if err != nil {
				return shim.Error(err.Error())
			}
		}
		resultBytes, err := stub.GetState(compositKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		result = string(resultBytes[:])
	} else if fcn == "write" {
		err := stub.PutState(compositKey, []byte(args[1]))
		if err != nil {
			return shim.Error(err.Error())
		}
		result = args[0] + "=>" + args[1]
	} else if fcn == "delete" {
		err := stub.DelState(compositKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		result = args[0]
	} else if fcn == "worldStates" {
		kvs, err:=WorldStates(stub);
		if err != nil {
			return shim.Error(err.Error())
		}
		var strings []string;
		for _, kv := range kvs {
			jsonElement:= KV2Json(kv)
			strings = append(strings,jsonElement)
		}
		jsonBytes,_:=json.Marshal(strings)
		result = string(jsonBytes)
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
