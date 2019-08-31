#!/bin/bash

for GOOS in darwin linux windows; do
    for GOARCH in amd64; do
        go build -v -o ./bin/rockit-$GOOS
    done
done