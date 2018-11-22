GO ?= go
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/)
VETPACKAGES ?= $(shell $(GO) list ./... | grep -v /vendor/ | grep -v /examples/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
TESTFOLDER := $(shell $(GO) list ./... | grep -E 'gin$$|binding$$|render$$' | grep -v examples)

all: install

install: deps
	govendor sync

.PHONY: test
test:
	echo "mode: count" > coverage.out
	for d in $(TESTFOLDER); do \
		$(GO) test -v -covermode=count -coverprofile=profile.out $$d; \
		if [ -f profile.out ]; then \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out; \
		fi; \
	done

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

vet:
	$(GO) vet $(VETPACKAGES)

deps:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/kardianos/govendor; \
	fi
	@hash embedmd > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/campoy/embedmd; \
	fi

embedmd:
	embedmd -d *.md

.PHONY: lint
lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

.PHONY: tools
tools:
	go install golang.org/x/lint/golint; \
	go install github.com/client9/misspell/cmd/misspell; \
	go install github.com/campoy/embedmd;
