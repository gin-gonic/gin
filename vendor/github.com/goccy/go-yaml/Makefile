## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
TESTMOD := testdata/go_test.mod

$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: test
test:
	go test -v -race ./...
	go test -v -race ./testdata -modfile=$(TESTMOD)

.PHONY: simple-test
simple-test:
	go test -v ./...
	go test -v ./testdata -modfile=$(TESTMOD)

.PHONY: fuzz
fuzz:
	go test -fuzz=Fuzz -fuzztime 60s

.PHONY: cover
cover:
	go test -coverpkg=.,./ast,./lexer,./parser,./printer,./scanner,./token -coverprofile=cover.out -modfile=$(TESTMOD) ./... ./testdata

.PHONY: cover-html
cover-html: cover
	go tool cover -html=cover.out

.PHONY: ycat/build
ycat/build: $(LOCALBIN)
	cd ./cmd/ycat && go build -o $(LOCALBIN)/ycat .

.PHONY: lint
lint: golangci-lint ## Run golangci-lint
	@$(GOLANGCI_LINT) run

.PHONY: fmt
fmt: golangci-lint ## Ensure consistent code style
	@go mod tidy
	@go fmt ./...
	@$(GOLANGCI_LINT) run --fix

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

## Tool Versions
GOLANGCI_VERSION := 2.1.2

.PHONY: golangci-lint
.PHONY: $(GOLANGCI_LINT)
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	@test -s $(LOCALBIN)/golangci-lint && $(LOCALBIN)/golangci-lint version --short | grep -q $(GOLANGCI_VERSION) || \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) v$(GOLANGCI_VERSION)
