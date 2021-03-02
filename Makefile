test:
	UNIT_TESTS=true ROOT_DIRECTORY=$(shell pwd) go test ./...

release:
	bash scripts/createrelease.sh
