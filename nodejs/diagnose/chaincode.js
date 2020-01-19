const {CommonChaincode, shim, ClientIdentity, format: {ParseHistory, ParseStates}} = require('fabric-common-chaincode');

const name = 'nodeDiagnose';

class Chaincode extends CommonChaincode {

	constructor() {
		super(name);
	}

	async init(stub) {
		return '';
	}

	/**
	 *
	 * @param {ChaincodeStub} stub
	 * @returns {Promise<string>}
	 */
	async invoke(stub) {
		const {fcn, params} = stub.getFunctionAndParameters();
		this.logger.info('Invoke', fcn);
		this.logger.debug('params', params);

		let response;
		switch (fcn) {
			case 'panic':
				throw Error('test panic');
			case 'transient':
				const key = params[0];
				const value = stub.getTransient(key);
				response = JSON.stringify({[key]: value});


				break;
			// case "richQuery":
			// 	const query = params[0];
			// 	this.Logger.Info("Query string", query)
			// 	const queryIter = this.stub.GetQueryResult(query)
			// 	const states = ParseStates(queryIter)
			// 	response = JSON.stringify(states)
			// case "worldStates":
			// 	var states = t.WorldStates("", nil)
			// 	response = ToJson(states)
			case 'whoami':
				response = (new ClientIdentity(stub.stub)).toString();
				break;
			case 'put':
				await stub.putState(params[0], params[1]);
				break;
			case 'get':
				response = await stub.getState(params[0]);
				break;
			default:
				this.logger.info('fcn:default');
				response = '';
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
			// case "putBatch":
			// 	var batch map[string]string
			// 	FromJson([]byte(params[0]), &batch)
			// 	for k, v := range batch {
			// 	t.PutState(k, []byte(v))
			// }
			// case "chaincodeId":
			// 	response = []byte(t.GetChaincodeID())
			// case "getCertID":
			// 	var certID = cid.NewClientIdentity(stub).GetID()
			// 	response = []byte(certID)
		}
		return response;
	}
}

shim.start(new Chaincode());
