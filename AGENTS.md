# AGENTS.md

Guidance for coding agents working on numbers.

## Project Overview

Microservice written in Go (`github.com/icco/numbers`) providing number manipulation, facts, and mathematical endpoints over HTTP.

## Commands

```sh
go test ./...    # Run tests
go vet ./...     # Vet code
go run main.go   # Run server locally (port 8080 by default)
go build .       # Build binary
```

## Architecture & Conventions

- `main.go` — Entrypoint, HTTP routing, and math handlers.
- Follow icco Go conventions (`github.com/icco/gutil` for logging and JSON rendering).
- PR titles and commits must follow Conventional Commits with lowercase subjects.
- Ensure all tests pass before submitting PRs.
