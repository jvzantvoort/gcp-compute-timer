NAME := gcp-compute-timer
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'
GOIMPORTS ?= goimports
GOCILINT ?= golangci-lint
GO ?= GO111MODULE=on go
.DEFAULT_GOAL := help

.PHONY: fmt
fmt: ## Formatting source codes.
	@$(GOIMPORTS) -w ./gcp-compute-timer

.PHONY: tags
tags:
	@find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*'| sed "s,$${PWD}/,," | xargs gotags >tags

.PHONY: pretty
pretty:
	@find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*' -exec goimports -w "{}" \;; \
	find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*' -exec gofmt -w "{}" \;

.PHONY: lint
lint: ## Run golint and go vet.
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell ./...

.PHONY: test
test:  ## Run the tests.
	@$(GO) test ./...

.PHONY: build
build: main.go  ## Build a binary.
	$(GO) build -ldflags "$(LDFLAGS)"

.PHONY: cross
cross: main.go  ## Build binaries for cross platform.
	mkdir -p pkg
	@# darwin
	@for arch in "amd64" "386"; do \
		GOOS=darwin GOARCH=$${arch} make build; \
		zip pkg/gcp-compute-timer_$(VERSION)_darwin_$${arch}.zip gcp-compute-timer; \
	done;
	@# linux
	@for arch in "amd64" "386" "arm64" "arm"; do \
		GOOS=linux GOARCH=$${arch} make build; \
		zip pkg/gcp-compute-timer_$(VERSION)_linux_$${arch}.zip gcp-compute-timer; \
	done;
