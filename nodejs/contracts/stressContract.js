const BaseContract = require('khala-fabric-contract-api/baseContract');


class StressContract extends BaseContract {

	constructor() {
		super('stress');
	}

	async init(context) {
		super.init(context);
	}

}

module.exports = StressContract;
