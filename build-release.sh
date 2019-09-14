#!/bin/bash

set -e

version=$(git tag)

# android & linux
GOOS=linux CGO_ENABLED=0 GOARCH=arm GOARM=7 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-linux-arm"
GOOS=linux CGO_ENABLED=0 GOARCH=arm64 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-linux-arm64"
GOOS=linux CGO_ENABLED=0 GOARCH=386 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-linux-i386"
GOOS=linux CGO_ENABLED=0 GOARCH=amd64 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-linux-amd64"

# darwin
GOOS=darwin CGO_ENABLED=0 GOARCH=386 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-darwin-i386"
GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-darwin-amd64"

# windows
GOOS=windows CGO_ENABLED=0 GOARCH=386 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-windows-i386.exe"
GOOS=windows CGO_ENABLED=0 GOARCH=amd64 \
    go build -gcflags "-trimpath $PWD" -ldflags "-s -w" -o "out/tlswrap_$version-windows-amd64.exe"

upx out/*
