set -e
CURRENT=$(cd $(dirname ${BASH_SOURCE}) && pwd)
CHAINCODE_SOURCE_DIR=${CURRENT}
CHAINCODE_METADATA_DIR=${CURRENT}/metadata
BUILD_OUTPUT_DIR=${CURRENT}/dst
RUN_METADATA_DIR=${CURRENT}/dst
detect() {
	bin/detect "$CHAINCODE_SOURCE_DIR" "$CHAINCODE_METADATA_DIR"
}
build() {
	bin/build "$CHAINCODE_SOURCE_DIR" "$CHAINCODE_METADATA_DIR" "$BUILD_OUTPUT_DIR"
}
release() {
	# bin/release script is responsible for providing the connection.json to the peer
	# by placing connection.json in the RELEASE_OUTPUT_DIR
	bin/release "$BUILD_OUTPUT_DIR" "$RELEASE_OUTPUT_DIR"

}
run() {
	bin/run "$BUILD_OUTPUT_DIR" "$RUN_METADATA_DIR"
}
up() {
	detect
	build
	release
	run
}
"$@"
