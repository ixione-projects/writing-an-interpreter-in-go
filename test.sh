#!/usr/bin/env bash

packages=""
for dir in "$(ls -d $HOME/src/writing-an-interpreter-in-go/src/go/*/)"; do
    packages="$packages $dir"
done

go test $packages
