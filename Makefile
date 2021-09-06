V=0
Q=$(if $(filter 1,$V),,@)
M=$(shell printf "\033[34;1m▶\033[0m")
SHELL:=/bin/bash

PRICECHECKER_FILES=go.mod go.sum $(shell find -type f -name "*.go")

GLOBAL_DOCKER_ARGS:=--rm --interactive --tty --volume $(PWD)/.cache:/tmp/cache --volume $(PWD):/workspace --workdir /workspace

GO_DOCKER_IMAGE:=golang:1.17
GOCACHE:=/tmp/cache/gocache
GOPATH:=/tmp/cache/gopath
GO_DOCKER_ARGS:=-e GOPATH=$(GOPATH) -e GOCACHE=$(GOCACHE)

GO_LINT_IMAGE:=golangci/golangci-lint:v1.42
GO_LINT_LINTERS:=-E stylecheck -E gofmt -E exhaustive -E makezero

### Targets

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Cleanup built artifacts
	$(info $(M) removing dist/...)
	$Q rm -rf dist/

.PHONY: test
test: test-pricechecker ## Run all project tests

.PHONY: lint
lint: lint-pricechecker

dist:
	$Q mkdir dist

.PHONY: build
build: dist/pricechecker ## Build binary

marketdata/items.go: marketdata/gen.go
	$(info $(M) running go generate...)
	$Q docker run \
		$(GLOBAL_DOCKER_ARGS) \
		$(GO_DOCKER_ARGS) \
		$(GO_DOCKER_IMAGE) \
		bash -c 'go generate'

.PHONY: test-pricechecker
test-pricechecker:
	$(info $(M) testing pricechecker...)
	$Q docker run \
		$(GLOBAL_DOCKER_ARGS) \
		$(GO_DOCKER_ARGS) \
		$(GO_DOCKER_IMAGE) \
		bash -c 'go test -v -race ./...'

.PHONY: lint-pricechecker
lint-pricechecker:
	$(info $(M) linting pricechecker...)
	$Q docker run \
		$(GLOBAL_DOCKER_ARGS) \
		$(GO_DOCKER_ARGS) \
		$(GO_LINT_IMAGE) \
		bash -c 'golangci-lint run $(GO_LINT_LINTERS) ./...'


dist/pricechecker: $(PRICECHECKER_FILES) marketdata/items.go | dist
	$(info $(M) building pricechecker...)
	$Q docker run \
		$(GLOBAL_DOCKER_ARGS) \
		$(GO_DOCKER_ARGS) \
		$(GO_DOCKER_IMAGE) \
		bash -c 'go build -o dist/pricechecker'
