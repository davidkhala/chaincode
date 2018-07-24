const {Base, shim} = require('../common/base');


class Chaincode extends Base {
	constructor() {
		const name = "test";
		super(name);
	}
}

shim.start(new Chaincode());
