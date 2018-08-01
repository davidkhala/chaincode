const {Base, shim, ChaincodeStub} = require('./common/base');


class Chaincode extends Base {
	constructor() {
		const name = "test";
		super(name);
	}

	async init(stub, clientIdentity) {
		return '';
	}

	/**
	 * @param {ChaincodeStub} stub
	 * @param {ClientIdentity} clientIdentity
	 * @returns {Promise<string>}
	 */
	async invoke(stub, clientIdentity) {

		this.setEvent("chaincodeEvent", "Hello World");
		const collection = 'private1';
		const key = '1';
		await stub.getPrivateData(collection, key);
		await stub.putPrivateData(collection, key, '');
		return '';
	}
}

shim.start(new Chaincode());
