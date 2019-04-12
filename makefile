TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=ovc
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: fmtcheck
	go build -mod vendor ${BUILDARGS}

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
	@golangci-lint run ./$(PKG_NAME)
	@go vet ./$(PKG_NAME)

test-compile: fmtcheck
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: default build install test fmt fmtcheck lint tools test-compile
