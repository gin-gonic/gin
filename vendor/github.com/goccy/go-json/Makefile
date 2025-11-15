PKG := github.com/goccy/go-json

BIN_DIR := $(CURDIR)/bin
PKGS := $(shell go list ./... | grep -v internal/cmd|grep -v test)
COVER_PKGS := $(foreach pkg,$(PKGS),$(subst $(PKG),.,$(pkg)))

COMMA := ,
EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
COVERPKG_OPT := $(subst $(SPACE),$(COMMA),$(COVER_PKGS))

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

.PHONY: cover
cover:
	go test -coverpkg=$(COVERPKG_OPT) -coverprofile=cover.out ./...

.PHONY: cover-html
cover-html: cover
	go tool cover -html=cover.out

.PHONY: lint
lint: golangci-lint
	$(BIN_DIR)/golangci-lint run

golangci-lint: | $(BIN_DIR)
	@{ \
		set -e; \
		GOLANGCI_LINT_TMP_DIR=$$(mktemp -d); \
		cd $$GOLANGCI_LINT_TMP_DIR; \
		go mod init tmp; \
		GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.48.0; \
		rm -rf $$GOLANGCI_LINT_TMP_DIR; \
	}

.PHONY: generate
generate:
	go generate ./internal/...
