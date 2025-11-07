#!/usr/bin/env bash

if [[ "$#" -lt 1 ]]; then
    echo 'Usage: test.sh <packages>'
    exit 1
fi

packages=""
for arg in "$@"; do
    packages="$packages $HOME/src/writing-an-interpreter-in-go/src/go/$arg"
done

go test $packages
