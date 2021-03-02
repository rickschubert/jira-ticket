CURR_DIR=$(shell pwd)

test:
	UNIT_TESTS=true ROOT_DIRECTORY=${CURR_DIR} go test ./...

release:
	bash scripts/createrelease.sh
