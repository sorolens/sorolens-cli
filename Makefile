BINARY      := sorolens
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS     := -ldflags "-s -w -X github.com/sorolens/sorolens-cli/cmd.Version=$(VERSION)"
GOFLAGS     ?=

.PHONY: build install test update-golden lint release-dry-run clean help

build: ## Compile the sorolens binary
	go build $(GOFLAGS) $(LDFLAGS) -o $(BINARY) .

install: ## Install sorolens binary to GOPATH/bin
	go install $(GOFLAGS) $(LDFLAGS) .

test: ## Run all tests
	go test $(GOFLAGS) ./...

update-golden: ## Regenerate golden test fixtures
	UPDATE_GOLDEN=1 go test ./internal/format/...

lint: ## Run go vet and golangci-lint
	go vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not found, skipping"

release-dry-run: ## Build all release targets locally without publishing
	goreleaser release --snapshot --clean

clean: ## Remove built binary and dist/
	rm -f $(BINARY)
	rm -rf dist/

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
