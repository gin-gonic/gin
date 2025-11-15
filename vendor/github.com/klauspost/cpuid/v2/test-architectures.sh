#!/bin/sh

set -e

go tool dist list | while IFS=/ read os arch; do
    echo "Checking $os/$arch..."
    echo " normal"
    GOARCH=$arch GOOS=$os go build -o /dev/null .
    echo " noasm"
    GOARCH=$arch GOOS=$os go build -tags noasm -o /dev/null .
    echo " appengine"
    GOARCH=$arch GOOS=$os go build -tags appengine -o /dev/null .
    echo " noasm,appengine"
    GOARCH=$arch GOOS=$os go build -tags 'appengine noasm' -o /dev/null .
done
