TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=ovc
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: fmtcheck
	go build -mod vendor ${BUILDARGS}

build-linux: fmtcheck
	GOOS=linux GOARCH=amd64 go build -mod vendor ${BUILDARGS} .
	
build-darwin: fmtcheck
	GOOS=darwin GOARCH=amd64 go build -mod vendor ${BUILDARGS} .

install: fmtcheck
	go install -mod vendor

test: lint
	go test $(TEST) -timeout=30s -parallel=4

fmtcheck:
	@echo "==> Checking that code complies with gofmt requirements...."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

tools:
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

lint: fmtcheck
	@echo "==> Checking source code against linters..."
	@# errcheck dissabled as it errors on not handling d.Set which would be combursome for the likeiness of occurance
	@# inproper err handling will still be caught by ineffassign and other linters
	@golangci-lint run ./$(PKG_NAME) -D errcheck
	@go vet ./$(PKG_NAME)

.PHONY: default build install test fmt fmtcheck lint tools
