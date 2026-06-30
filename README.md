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

- Read-only virtual filesystem backed by embedded plain-text files
- Familiar commands including `ls`, `cd`, `cat`, `tree`, and `grep`
- Multi-stage pipelines
- Session-specific working directories and command history
- Command and path autocomplete
- HTTP API and disposable browser terminal
- Docker packaging and GitHub Actions CI

## Architecture

```text
Browser terminal
      ↓
HTTP API
      ↓
Shell engine
      ↓
Parser → Commands
             ↓
      Virtual filesystem
             ↓
    Embedded text files
```

Commands interact only with the virtual filesystem. Content is embedded into
the binary from `content/` at build time. Vaultsh exposes no commands that
create, modify, or delete files.

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

## Content

Portfolio content is stored as grep-friendly plain text:

```text
content/
├── about.txt
├── skills.txt
├── interests.txt
├── experience/
│   ├── reversinglabs.txt
│   ├── intellexi.txt
│   ├── a1.txt
│   └── arisglobal.txt
└── projects/
    └── vaultsh.txt
```

Each non-empty line uses a lowercase `key: value` format. Keys may repeat and
blank lines may separate sections.

```text
focus: backend services
focus: distributed systems

technology: Go
technology: Docker
```

Run `go test ./...` after editing content. The test suite validates the
embedded layout and file format.

## Project Structure

```text
cmd/vaultsh/          application entry point
content/              embedded portfolio content
internal/command/     shell commands
internal/filesystem/  read-only virtual filesystem
internal/httpapi/     HTTP transport
internal/parser/      tokenizer, lexer, and parser
internal/shell/       execution engine and sessions
internal/storage/     embedded-content loader
testui/               browser terminal
```

## Documentation

- [Command reference and examples](docs/commands.md)
- [Content layout and format](docs/content.md)
- [HTTP API](docs/api.md)
- [Roadmap](docs/roadmap.md)
