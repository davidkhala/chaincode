const {Contract} = require('fabric-contract-api');


class StressContract extends Contract {

	constructor() {
		super('stress');
	}

	async invoke(context) {
		this.logger.info('invoke', this.getName());
	}

}

module.exports = StressContract;
