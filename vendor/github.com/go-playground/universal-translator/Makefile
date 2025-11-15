GOCMD=GO111MODULE=on go

linters-install:
	@golangci-lint --version >/dev/null 2>&1 || { \
		echo "installing linting tools..."; \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.41.1; \
	}

lint: linters-install
	golangci-lint run

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -bench=. -benchmem ./...

.PHONY: test lint linters-install