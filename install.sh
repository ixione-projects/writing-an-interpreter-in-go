#!/usr/bin/env bash

go build -o monkey $HOME/src/go/ && mkdir -p $HOME/go/bin/ && mv monkey $HOME/go/bin/
