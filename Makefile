
test:
	@go test ./... -v

test-ci:
	@go test ./... 

.PHONY: test test-ci
