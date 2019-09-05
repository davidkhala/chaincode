# Chaincode


## Notes
- call `await stub.putPrivateData('any', "key", 'value');` without setup collection Config or in Init step:  
Error: collection config not define for namespace [node]
See also in https://github.com/hyperledger/fabric/commit/8a705b75070b7a7021ec6f897a80898abf6a1e45
- random function should not be used in chaincode as part of ledger data, otherwise it corrupts the deterministic process.
- instantiate/upgrade could be where data migration is performed, if necessary

## Version Map
- 1.4.2: golang: 1.11.x


### TODO

