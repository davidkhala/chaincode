const BaseContract = require('khala-fabric-contract-api/baseContract');


class StressContract extends BaseContract {

	constructor() {
		super('stress');
	}

	async init(context) {
		super.init(context);
	}

	async panic(context) {
		throw Error('panic');
	}

	async unknownTransaction(context) {
		super.unknownTransaction(context);
		throw Error('unknownTransaction');
	}
}

module.exports = StressContract;
