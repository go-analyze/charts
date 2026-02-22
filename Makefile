export GO111MODULE = on

.PHONY: default test test-cover bench lint

# Packages to test, exclude packages without tests to avoid example noise
CODE_PKGS := $(shell go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./...)


test:
	go test -race -cover $(CODE_PKGS)

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench:
	$(eval CORES_HALF := $(shell expr `getconf _NPROCESSORS_ONLN` / 2))
	go test -parallel=$(CORES_HALF) --benchmem -benchtime=20s -bench='Benchmark.*Render' -run='^$$'

lint:
	golangci-lint run --timeout=600s
	go vet ./...
