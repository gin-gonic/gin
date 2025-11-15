#!/usr/bin/env bash

current=$(git status | head -n1 | sed 's/On branch //')
name=${1:-$current}
if [[ ! $name =~ ^(((opt(imize)?|feat(ure)?|(bug|hot)?fix|test|refact(or)?|ci)/.+)|(main|develop)|(release-v[0-9]+\.[0-9]+)|(release/v[0-9]+\.[0-9]+\.[0-9]+(-[a-z0-9.]+(\+[a-z0-9.]+)?)?)|revert-[a-z0-9]+)$ ]]; then
    echo "branch name '$name' is invalid"
    exit 1
else
    echo "branch name '$name' is valid"
fi
