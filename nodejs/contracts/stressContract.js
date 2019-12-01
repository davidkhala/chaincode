const {Contract, Context} = require('fabric-contract-api');


class StressContract extends Contract {

	constructor() {
		super('-');
	}

	async invoke(context) {
		console.log('overlapping invoke');
	}

	async init(context) {
		console.log('overlapping init');
	}
}

module.exports = StressContract;
