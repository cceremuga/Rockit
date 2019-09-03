#!/bin/bash

for GOOS in darwin linux windows freebsd openbsd; do
    for GOARCH in 386 amd64 arm arm64; do
        go build -v -o bin/rockit-$GOOS-$GOARCH
    done
done