#!/usr/bin/env sh

# Deal with go stuff...
go mod download golang.org/x/mobile

# To support macOS, use this fork of gomobile: https://github.com/ydnar/gomobile/tree/support-catalyst Download code, go run ./cmd/gobind and go run ./cmd/gomobile, then gomobile init
# Build
GO11MODULE=on gomobile bind -target=ios ./pkg/appnet

echo Done
