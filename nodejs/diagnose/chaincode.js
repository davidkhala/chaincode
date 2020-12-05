const {CommonChaincode, shim, ClientIdentity, format: {ParseHistory, ParseStates}} = require('fabric-common-chaincode');

const name = 'nodeDiagnose';

class Chaincode extends CommonChaincode {

	constructor() {
		super(name);
	}

	async init() {
		return '';
	}

	async panic() {
		throw Error('test panic');
	}

	async transient(key) {
		const value = this.getTransient(key);
		return {[key]: value};
	}

	// TODO
	async richQuery(query) {
		this.logger.info('Query string', query);
		const queryIter = this.stub.getQueryResult(query);
		const states = await ParseStates(queryIter);
		return states;
	}

	async putBatch(batch) {
		for (const [key, value] of Object.entries(JSON.parse(batch))) {
			await this.putState(key, value);
		}
	}

	async worldStates() {
		const iterator = await this.stub.getStateByRange();
		return await ParseStates(iterator);
	}

	async whoami() {
		return new ClientIdentity(this.stub).toString();
	}

	async put(key, value) {
		await this.putState(key, value);
	}

	async get(key) {
		return await this.stub.getState(key);
	}

	async history(key) {
		const iterator = await this.stub.getHistoryForKey(key);
		const history = await ParseHistory(iterator);
		return history;
	}

	async timeStamp() {
		return this.stub.getTxTimestamp(true);
	}

	// TODO
	// case "putEndorsement":
	// 	var key = params[0]
	// 	var orgs = params[1:]
	// 	var policy = ext.NewKeyEndorsementPolicy(nil)
	// 	policy.AddOrgs(msp.MSPRole_MEMBER, orgs...)
	// 	t.SetStateValidationParameter(key, policy.Policy())
	// case "getEndorsement":
	// 	var key = params[0]
	// 	var policyBytes = t.GetStateValidationParameter(key)
	// 	var policy = ext.NewKeyEndorsementPolicy(policyBytes)
	// 	var orgs = policy.ListOrgs()
	// 	response = ToJson(orgs)
	// case "delegate":
	// 	type crossChaincode struct {
	// 	ChaincodeName string
	// 	Fcn           string
	// 	Args          []string
	// 	Channel       string
	// }
	// 	var paramInput crossChaincode
	// 	FromJson([]byte(params[0]), &paramInput)
	// 	var args = ArgsBuilder(paramInput.Fcn)
	// 	for i, element := range paramInput.Args {
	// 	args.AppendArg(element)
	// 	t.Logger.Debug("delegated Arg", i, element)
	// }
	// 	var pb = t.InvokeChaincode(paramInput.ChaincodeName, args.Get(), paramInput.Channel)
	// 	response = pb.Payload
	// case "listPage":
	// 	var startKey = params[0]
	// 	var endKey = params[1]
	// 	var pageSize = Atoi(params[2])
	// 	var bookMark = params[3]
	// 	var iter, metaData = t.GetStateByRangeWithPagination(startKey, endKey, pageSize, bookMark)
	//
	// 	type Response struct {
	// 	States   []StateKV
	// 	MetaData QueryResponseMetadata
	// }
	// 	response = ToJson(Response{ParseStates(iter, nil), metaData})
	// case "chaincodeId":
	// 	response = []byte(t.GetChaincodeID())

}

shim.start(new Chaincode());
