const {Shim} = require('fabric-shim');

class Chaincode {

	constructor() {
	}

	async Init(stub) {
		return Shim.success();
	}

	async Invoke(stub) {
		return Shim.success();
	}
}

Shim.start(new Chaincode());
