const {CommonChaincode, shim} = require('fabric-common-chaincode');

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
		const value = await this.getPrivateData(collection, key);
		if (!value) {
			this.logger.info('clientID' + clientIdentity.getID());
			await this.putPrivateData(collection, key, clientIdentity.getID());
		}
		this.logger.info(`PrivateData done valueType ${typeof value},length:${value.length}, value: ${value}`);
		return '';
	}
}

shim.start(new Chaincode());
