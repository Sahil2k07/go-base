
.DEFAULT_GOAL := build

.PHONY: build server migration tests build-windows build-macos migration-windows migration-macos

# Default: Linux (amd64)
build: server

server:
	go build -o server ./cmd/server

# Windows (amd64)
server-windows:
	go build -o server.exe ./cmd/server

# macOS (amd64)
server-macos:
	go build -o server ./cmd/server

# Migration binaries
migration:
	go build -o migration ./cmd/migration

migration-windows:
	go build -o migration.exe ./cmd/migration

migration-macos:
	go build -o migration ./cmd/migration

gen:
	go run cmd/gen/main.go

# Tests
tests:
	go test -count=1 ./...
