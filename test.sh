#!/usr/bin/env bash

packages=""

if [[ $# -lt 1 ]]; then
    for dir in "$(ls -d $HOME/src/writing-an-interpreter-in-go/src/go/*/)"; do
        packages="$packages $dir"
    done
else
    for arg in "$@"; do
        packages="$packages $HOME/src/writing-an-interpreter-in-go/src/go/$arg"
    done
fi

go test $packages
