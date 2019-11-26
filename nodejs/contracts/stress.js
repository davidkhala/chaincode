const {Contract, Context} = require('fabric-contract-api');

class StressContract extends Contract {

	constructor() {
		// Supplying a name within the constructor is also option and will default to ''
		// name could be got by getName()
		super('stressContract');
	}

	/**
	 * This permits contracts to use their own subclass of context
	 * After this function returns, the chaincodeStub, client identity objects and logger will be injected.
	 * @override
	 * @returns {Context}
	 */
	createContext() {
		this.context = super.createContext();
		return this.context;
	}

	/**
	 *
	 * @param {string} [label]
	 */
	createLogger(label) {
		this.logger = this.context.logging.getLogger(label); // contract name has been added as tag before label
	}

	/**
	 * @override
	 * @param {Context} context
	 * @returns {Promise<void>}
	 */
	async beforeTransaction(context) {
		this.createLogger();
		await super.beforeTransaction(context);
	}

	/**
	 * @override
	 * @param {Context} context
	 * @param result
	 * @return {Promise<void>}
	 */
	async afterTransaction(context, result) {
		await super.afterTransaction(context, result);
	}

	/**
	 * @override
	 * @param {Context} context
	 * @return {Promise<void>}
	 */
	async unknownTransaction(context) {
		await super.unknownTransaction(context);
	}

	/**
	 * @override
	 * @return {string}
	 */
	getName() {
		return super.getName();
	}

	/**
	 * 'init' as a default fcn in sdk-node
	 * @param {Context} context the transaction context
	 */
	async init(context) {
		console.log(context);
	}
}

module.exports = StressContract;
