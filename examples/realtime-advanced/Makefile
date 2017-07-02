all: deps build

.PHONY: deps
deps:
	go get -d -v github.com/dustin/go-broadcast/...
	go get -d -v github.com/manucorporat/stats/...

.PHONY: build
build: deps
	go build -o realtime-advanced main.go rooms.go routes.go stats.go
