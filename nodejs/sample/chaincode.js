const {Base, shim} = require('./common/base');


class Chaincode extends Base {
	constructor() {
		const name = "test";
		super(name);
	}

	async init(stub, clientIdentity) {
		return '';
	}

	async invoke(stub, clientIdentity) {

		this.setEvent("chaincodeEvent", "Hello World");

		return '';
	}
}

shim.start(new Chaincode());
