#!/usr/bin/env sh

# Deal with go stuff...
go mod download golang.org/x/mobile

# Build
GO11MODULE=on gomobile bind -target=ios ./pkg/appnet

echo Done
