const { Contract, Context } = require('fabric-contract-api');

class StressContext extends Context {

	constructor() {
		super();
	}

}

class StressContract extends Contract {

	constructor() {
		// Unique namespace when multiple contracts per chaincode file
		super('stressContract');
	}

	createContext() {
		return new StressContext();
	}

	/**
	 * Instantiate to perform any setup of the ledger that might be required.
	 * @param {Context} ctx the transaction context
	 */
	async instantiate(ctx) {
		// No implementation required with this example
		// It could be where data migration is performed, if necessary

	}
}

module.exports = StressContract;
