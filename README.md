# vaultsh

A read-only virtual shell engine with a virtual filesystem, sessions, pipelines,
and pluggable commands.

## Why

Vaultsh started as a backend-first portfolio experiment.

Instead of building a traditional static portfolio, the idea was to expose
curated profile content through a small shell-like system. The project provides
room to explore command parsing, virtual filesystems, stateful sessions,
pipelines, HTTP APIs, testing, and container deployment.

## Features

- Read-only virtual filesystem with embedded content
- Familiar commands including `ls`, `cd`, `cat`, `tree`, and `grep`
- Multi-stage pipelines
- Session-specific working directories and command history
- Command and path autocomplete
- HTTP API and disposable browser terminal
- Docker packaging and GitHub Actions CI

## Quickstart

Prerequisites:

- Docker Engine or Docker Desktop
- Docker Compose v2

```sh
docker compose up --build
```

Open:

```text
http://localhost:8080/testui/
```

Stop:

```sh
docker compose down
```

For local development, install Go 1.24 or newer:

```sh
go run ./cmd/vaultsh
go test ./...
```

## Example Commands

```sh
help
tree
ls -la /
cat about.txt
cat skills.txt | grep "^language:" | sort
history | tail -n 5
```

## Documentation

- [Command reference and examples](docs/commands.md)
- [Content layout and format](docs/content.md)
- [HTTP API](docs/api.md)
- [Roadmap](docs/roadmap.md)
