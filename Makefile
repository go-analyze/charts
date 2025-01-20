export GO111MODULE = on

.PHONY: default test test-cover bench lint


test:
	go test -race -cover ./...

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench:
	$(eval CORES_HALF := $(shell expr `getconf _NPROCESSORS_ONLN` / 2))
	go test -parallel=$(CORES_HALF) --benchmem -benchtime=20s -bench=. -run='^$$' ./...

lint:
	golangci-lint run
