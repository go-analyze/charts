export GO111MODULE = on

.PHONY: default test test-cover bench lint


test:
	go test -race -cover ./...

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench:
	go test --benchmem -bench=. ./...

lint:
	golangci-lint run
