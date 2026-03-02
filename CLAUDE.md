# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`sqlight` is a Go module (`github.com/netomathias/sqlight`) targeting Go 1.25.

## Common Commands

```bash
# Build
go build ./...

# Run all tests
go test ./...

# Run a single test
go test ./... -run TestName

# Run tests with verbose output
go test -v ./...

# Lint (requires golangci-lint)
golangci-lint run
```
