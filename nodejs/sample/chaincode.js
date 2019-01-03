const {CommonChaincode, shim, ClientIdentity} = require('fabric-common-chaincode');

const collection = 'private1';
const key = '1';

class Chaincode extends CommonChaincode {

	constructor() {
		const name = 'test';
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
		this.setEvent('chaincodeEvent', 'Hello World');
		let value = await this.getPrivateData(collection, key);
		if (!value) {
			const clientIdentity = new ClientIdentity(stub);
			value = clientIdentity.getID();
			await this.putPrivateData(collection, key, value);
		}
		return value;
	}
}

shim.start(new Chaincode());
