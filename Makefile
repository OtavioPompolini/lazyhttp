.DEFAULT_GOAL:=help
SHELL:=/bin/sh

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

########################################################################################################################
##@ Build
########################################################################################################################

.PHONY: build
build: ## Builds the app
	go build -ldflags="-w -s" -v -o ./bin/lazyhttp ./cmd/lazyhttp

########################################################################################################################
##@ Test
########################################################################################################################

.PHONY: test
test: ## Runs all the tests
	go test -v ./...

.PHONY: test-cover
test-cover: ## Runs all the tests with coverage
	go test -v -cover ./...

########################################################################################################################
##@ Run
########################################################################################################################

.PHONY: run
run: ## Runs the app
	go run ./cmd/lazyhttp

########################################################################################################################
##@ Code Style
########################################################################################################################

.PHONY: format
format: FORMAT_PATH:=.
format: ## Runs goimports to format code and organize imports on the specified path (defaults to .)
	goimports --local "github.com/OtavioPompolini/project-postman" -w $(FORMAT_PATH)

.PHONY: lint
lint: LINT_OPTIONS:=
lint: ## Runs golangci-lint
	golangci-lint run $(LINT_OPTIONS)

.PHONY: lint-fix
lint-fix: ## Runs golangci-lint with the --fix flag
	golangci-lint run --fix

########################################################################################################################
##@ Tools
########################################################################################################################

.PHONY: generate
generate: ## Runs go generate accross all packages
	go generate ./...

.PHONY: clean-mocks
clean-mocks: ## Remove all the application's mocks
	find . -name 'mock_*.go' -delete
