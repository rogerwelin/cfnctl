# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

cfnctl is a Go CLI tool that brings a Terraform-like workflow to AWS CloudFormation. Commands: `apply`, `plan`, `destroy`, `output`, `validate`, `version`. Built with urfave/cli v2 and AWS SDK v2.

## Build & Development Commands

```bash
make build          # Build binary (CGO_ENABLED=0 go build ./cmd/cfnctl)
make test           # Run tests with coverage + race detection (go test -cover -race -v ./...)
make lint           # Run golangci-lint (config in .golangci.yml)
make check          # Run test + lint
make audit          # Tidy, verify deps, lint
```

Run a single test:
```bash
go test -run TestFunctionName -v ./commands/
```

## Architecture

**Entry point**: `cmd/cfnctl/main.go` → `cli/cli.go` (command registration)

**Command flow**: CLI flags parsed in `cli/` → command structs (`cli/types.go`) call `Run()` → implementations in `commands/` → CloudFormation API calls via `pkg/client/`

Key layers:
- **`cli/`** — CLI setup, flag definitions, parameter handling (`cli/params/` for YAML param files and interactive prompts)
- **`commands/`** — Command logic. Each command uses `CommandBuilder` (`helper.go`) to create a configured `Cfnctl` client, then calls methods on it
- **`pkg/client/`** — `Cfnctl` struct wraps CloudFormation operations. `CloudformationAPI` interface enables mock testing. Uses functional options pattern for configuration
- **`aws/`** — AWS SDK client initialization
- **`internal/interactive/`** — Real-time stack event streaming with color-coded table output
- **`internal/mock/`** — Mock implementation of `CloudformationAPI` interface for tests

**CloudFormation workflow in apply**: CreateChangeSet → DescribeChangeSet → (user approval) → ExecuteChangeSet → poll DescribeStack until terminal state, streaming resource events to table

## Testing

Tests live alongside commands in `commands/*_test.go`. They use the mock service in `internal/mock/mocksvc.go` which implements the `CloudformationAPI` interface from `pkg/client/types.go`. Test fixtures are in `commands/testdata/`.

## CI

GitHub Actions (`.github/workflows/ci.yml`): lint job runs golangci-lint, test job runs `make test` on stable + oldstable Go. GoReleaser handles cross-platform release builds.
