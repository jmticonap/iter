.PHONY: test

test:
	@go test -v -cover ./test/... -coverpkg=./...