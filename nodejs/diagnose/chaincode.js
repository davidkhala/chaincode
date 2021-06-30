const {CommonChaincode, shim, ClientIdentity, format: {parseHistory, parseStates}} = require('fabric-common-chaincode');

/**
 *
 * @typedef {Object} TxData
 * @property {number} Time
 * @property {byte[]} Value
 */


class Chaincode extends CommonChaincode {

	constructor() {
		super();
	}

	async init() {
		return '';
	}

	async panic() {
		throw Error('test panic');
	}

	async transient(key) {
		const value = this.stub.getTransient(key);
		return {[key]: value};
	}

	async richQuery(query) {
		this.logger.info('Query string', query);
		const queryIterator = await this.stub.getQueryResult(query);
		return await parseStates(queryIterator);
	}

	async putBatch(batch) {
		for (const [key, value] of Object.entries(JSON.parse(batch))) {
			await this.putRaw(key, value);
		}
	}

	async worldStates() {
		const iterator = await this.stub.getStateByRange();
		return await parseStates(iterator);
	}

	async whoami() {
		const cid = new ClientIdentity(this.stub.stub);
		return cid.toString();
	}

	async putRaw(key, value) {
		await this.stub.putState(key, value);
	}

	async getRaw(key) {
		return await this.stub.getState(key);
	}

	async put(key, value) {
		const data = {
			Time: this.timeStamp(), Value: value
		};
		await this.putRaw(key, JSON.stringify(data));
	}

	async history(key) {
		const iterator = await this.stub.getHistoryForKey(key);
		return await parseHistory(iterator);
	}

	timeStamp() {
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
	async chaincodeId() {
		return await this.stub.getChaincodeID();

	}

}

shim.start(new Chaincode());
