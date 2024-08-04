BIN_DIR := ./bin
GOLANGCI_LINT_VERSION := v1.59.1

.PHONY: install-linter
install-linter:
	wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LINT_VERSION) -b $(BIN_DIR)

.PHONY: lint
lint:
	$(BIN_DIR)/golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	$(BIN_DIR)/golangci-lint run ./... --fix

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: build
build:
	go build -o ./bin/ .
