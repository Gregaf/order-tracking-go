BIN_DIR := bin
SOURCE=$(shell find . -name "*.go" -not -path "*/infra/*")
TEST_DIRS=$(shell find . -name "*_test.go" -not -path "*/infra/*" -printf "%h\n" | sort -u)


INFRA_SOURCE=$(shell find ./infra -name "*.go")
CDK_OUT := infra/cdk.out/InfraStack.template.json

ENTRY_DIR := cmd
ENTRY_POINTS=$(shell find $(ENTRY_DIR) -iname "*.go")
BIN_PATHS := $(patsubst $(ENTRY_DIR)/%,$(BIN_DIR)/%/bootstrap,$(shell dirname $(ENTRY_POINTS)))

.DEFAULT_GOAL := help

lint: ## Lint the project
	golangci-lint run -v

test: $(TEST_DIRS) ## Run all tests

$(TEST_DIRS):
	cd $@ && go test -v

test-coverage: ## Run all tests and output coverage file
	@echo "TODO: Must implement makefile test coverage rule"	

build: $(BIN_PATHS) ## Build all entrypoints

$(BIN_DIR)/%/bootstrap: $(ENTRY_DIR)/% $(SOURCE)
	@echo "Building $@ from $</main.go"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@ $</main.go	

serve: ## Start the SAM Local API
	sam local start-api -t $(CDK_OUT) --warm-containers eager

synth: build $(CDK_OUT) ## Synthesize the CDK stack

$(CDK_OUT): $(INFRA_SOURCE)
	cd infra && cdk synth -v --no-staging

clean: ## Remove built artifacts
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)

help:  ## Display this help
	@$(info Order Tracking Application)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: help clean synth serve test lint test-coverage build