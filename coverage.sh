#!/usr/bin/env bash

set -e

echo "mode: count" > coverage.out

for d in $(go list ./... | grep -E 'gin$|binding$|render$' | grep -v 'examples'); do
    go test -v -covermode=count -coverprofile=profile.out $d
    if [ -f profile.out ]; then
        cat profile.out | grep -v "mode:" >> coverage.out
        rm profile.out
    fi
done
