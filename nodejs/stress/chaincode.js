const {Shim} = require('fabric-shim');

class Chaincode {

	constructor() {
	}

	async Init(stub) {
		// Init function should be implemented, even with init_required=false
		return Shim.success();
	}

	async Invoke(stub) {
		return Shim.success();
	}
}

Shim.start(new Chaincode());
