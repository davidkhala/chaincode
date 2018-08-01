const { Base, shim } = require('./common/base');

const collection = 'private1';
const key = '1';
class Chaincode extends Base {

	constructor() {
		const name = 'test';
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
