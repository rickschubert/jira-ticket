test:
	UNIT_TESTS=true go test ./...

release:
	bash scripts/createrelease.sh
