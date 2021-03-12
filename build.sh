#!/usr/bin/env sh

# Deal with go stuff...
go mod download golang.org/x/mobile

# To support macOS, use this fork of gomobile: https://github.com/ydnar/gomobile/tree/support-catalyst Download code, go run ./cmd/gobind and run: go build ./cmd/gobind && go build ./cmd/gomobile && cp gobind $GOPATH/bin/ && cp gomobile $GOPATH/bin/
# Build
GO11MODULE=on gomobile bind -target=ios ./pkg/appnet

echo Done
