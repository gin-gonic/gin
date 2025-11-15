SHELL := /bin/bash
RAGEL := ragel
GOFMT := go fmt

export GO_TEST=env GOTRACEBACK=all go test $(GO_ARGS)

.PHONY: build
build: machine.go

.PHONY: clean
clean:
	@rm -rf docs
	@rm -f machine.go

.PHONY: images
images: docs/urn.png

.PHONY: snake2camel
snake2camel:
	@cd ./tools/snake2camel; go build -o ../../snake2camel .

.PHONY: removecomments
removecomments:
	@cd ./tools/removecomments; go build -o ../../removecomments .

machine.go: machine.go.rl

machine.go: snake2camel

machine.go: removecomments

machine.go:
	$(RAGEL) -Z -G1 -e -o $@ $<
	@./removecomments $@
	@./snake2camel $@
	$(GOFMT) $@

docs/urn.dot: machine.go.rl
	@mkdir -p docs
	$(RAGEL) -Z -e -Vp $< -o $@

docs/urn.png: docs/urn.dot
	dot $< -Tpng -o $@

.PHONY: bench
bench: *_test.go machine.go
	go test -bench=. -benchmem -benchtime=5s ./...

.PHONY: tests
tests: *_test.go
	$(GO_TEST) ./...
