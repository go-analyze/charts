export GO111MODULE = on

.PHONY: default test test-cover bench lint

# Packages to test (exclude example packages)
CODE_PKGS := $(shell go list ./... | grep -v '/examples')


test:
	go test -race -cover $(CODE_PKGS)

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench:
	$(eval CORES_HALF := $(shell expr `getconf _NPROCESSORS_ONLN` / 2))
	go test -parallel=$(CORES_HALF) --benchmem -benchtime=20s -bench='Benchmark.*Render' -run='^$$'

lint:
	golangci-lint run --timeout=600s --enable=asasalint,asciicheck,bidichk,containedctx,contextcheck,decorder,durationcheck,errorlint,exptostd,fatcontext,forbidigo,gocheckcompilerdirectives,gochecksumtype,goconst,gofmt,goimports,gosmopolitan,grouper,iface,importas,mirror,misspell,nilerr,nilnil,perfsprint,prealloc,reassign,recvcheck,sloglint,testifylint,unconvert,wastedassign,whitespace && go vet ./...

