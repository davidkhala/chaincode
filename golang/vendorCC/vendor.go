package main

import (
	"encoding/json"
	. "github.com/davidkhala/fabric-common-chaincode-golang"
	"github.com/davidkhala/fabric-common-chaincode-golang/cid"
	. "github.com/davidkhala/goutils"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	name      = "vendor"
	TODO      = "TODO"
	Submitted = "Submitted"
	Confirmed = "Confirmed"
	Reject    = "Reject"
	Closed    = "Closed"
	StepType  = "Schedule"
)

type SimpleChaincode struct {
	CommonChaincode
}

func (t SimpleChaincode) Init(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Init", fcn)
	t.Logger.Debug("params", params)
	var project Project
	FromJson([]byte(params[0]), &project)
	var projectTitle = project.Title
	t.PutState(projectTitle, []byte(params[0]))

	for _, step := range project.Schedule {
		var key = t.CreateCompositeKey(StepType, []string{step.ID})
		t.PutState(key, ToJson(step))
	}
	return shim.Success(nil)
}

func (t SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) (response peer.Response) {
	defer Deferred(DeferHandlerPeerResponse, &response)
	t.Prepare(stub)

	var responseBytes []byte
	var fcn, params = stub.GetFunctionAndParameters()
	t.Logger.Info("Invoke", fcn)
	t.Logger.Debug("params", params)

	switch fcn {
	case "read":
		responseBytes = t.read(params[0], params[1])
	case "progress":
		t.progress(params)
	default:
		PanicString("Unknown fcn:" + fcn)
	}
	return shim.Success(responseBytes)
}

func (t SimpleChaincode) progress(params []string){


	var client = cid.NewClientIdentity(t.CCAPI)
	var step Step
	switch client.Cert.Subject.CommonName{
	case partBSubject:
		var submit Submit
		FromJson([]byte(params[0]),&submit)

		t.GetStateObj(submit.ID,&step)

		if step.Status != TODO && step.Status != Reject && step.Status != Submitted {
			PanicString("Invalid current step. Status:" + step.Status)
		}
		step.lastSubmit = submit

		step.Status = Submitted

	case partASubject:
		var review Review
		FromJson([]byte(params[0]),&review)

		t.GetStateObj(review.ID,&step)

		if step.Status != Submitted {
			PanicString("Invalid current step. Status:" + step.Status)
		}
		step.lastReview = review

		step.Status = review.Status

	case partCSubject:
		var audit Audit
		FromJson([]byte(params[0]),&audit)

		t.GetStateObj(audit.ID,&step)


		if step.Status != Confirmed {
			PanicString("Invalid current step. Status:" + step.Status)
		}
		step.lastAudit = audit

		step.Status = audit.Status

	default:
		PanicString("invalid creator:" + cert.Subject.CommonName)
	}

	t.PutStateObj(step.ID,step)

}

func (t SimpleChaincode) read(targetType, item string) []byte {
	var value []byte
	switch targetType {
	case "project":
		value = t.GetState(item)
	case "step":
		var key = t.CreateCompositeKey(StepType, []string{item})
		value = t.GetState(key)
	default:
		PanicString("Unknown read targetType:" + targetType)
	}
	return value
}

func main() {
	var cc = SimpleChaincode{}
	cc.SetLogger(name)
	ChaincodeStart(cc)

}
