const {Contract, Context} = require('fabric-contract-api');

class BaseContract extends Contract {

	constructor(name = '-') {
		// Supplying a name is required, otherwize init function could not be discovered
		// name could be got by getName()
		super(name);
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
		console.log(`creating logger[label=${label}]`);
		this.logger = this.context.logging.getLogger(label); // contract name has been added as tag before label
	}

	/**
	 * @override
	 * @param {Context} context
	 * @returns {Promise<void>}
	 */
	async beforeTransaction(context) {
		if (!this.logger) {
			this.createLogger();
		}
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
		this.logger.info(context);
	}

	/**
	 * 'invoke' as a default fcn in sdk-node
	 * @param {Context} context the transaction context
	 */
	async invoke(context) {
		this.logger.info('invoke', this.getName());
	}
}

module.exports = BaseContract;
