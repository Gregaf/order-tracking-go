ROOT := $(shell pwd)
DIST_DIR := $(ROOT)/dist
SOURCE := $(shell find $(ROOT) -name "*.go")

TEST_DIRS := $(shell find . -name "*_test.go" | xargs -I {} dirname {} | sort -u)
COVERAGE_FILE := coverage.out

ENTRY_DIR := $(ROOT)/app/cmd
ENTRY_POINTS := $(shell find $(ENTRY_DIR) -name "*.go")

BIN_DIR := $(ROOT)/bin
BIN_NAME := bootstrap
BIN_PATHS := $(patsubst $(ENTRY_DIR)/%,$(BIN_DIR)/%/$(BIN_NAME),$(shell dirname $(ENTRY_POINTS)))

.DEFAULT_GOAL := help

lint: ## Lint the project
	golangci-lint run -v

test: ## Run all tests
	go test $(TEST_DIRS) -v

test-cov: ## Run all tests and output coverage file
	go test $(TEST_DIRS) -coverprofile=$(COVERAGE_FILE) -v
	go tool cover -func=$(COVERAGE_FILE)	
	go tool cover -html=$(COVERAGE_FILE)

build: $(BIN_PATHS) ## Build all entrypoints

$(BIN_DIR)/%/$(BIN_NAME): $(ENTRY_DIR)/% $(SOURCE)
	@echo "Building $@ from $</main.go"
	cd $(ROOT)/app && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $@ $</main.go	

tidy:
	cd $(ROOT)/app && go mod tidy

clean: ## Remove built artifacts
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)

PACKAGES := $(patsubst $(BIN_DIR)/%/$(BIN_NAME),$(DIST_DIR)/%.zip,$(BIN_PATHS))

package: $(BUILD) $(DIST_DIR) $(PACKAGES) ## Create zip packages for all built binaries

$(DIST_DIR)/%.zip: $(BIN_DIR)/%/$(BIN_NAME)
	cd $(dir $<) && zip -r $(DIST_DIR)/$*.zip $(notdir $<)

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

help:  ## Display this help
	@$(info APPLICATION NAME)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: help clean lint test test-cov build tidy package
