#!/bin/bash

for PLATFORM in darwin linux windows; do
    for ARCH in 386 amd64; do
        export GOOS=$PLATFORM
        export GOARCH=$ARCH
        go build -o bin/rockit-$GOOS-$GOARCH
    done
done