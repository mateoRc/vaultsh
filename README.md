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

- Read-only virtual filesystem backed by mounted plain-text files
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
    Mounted text files
```

Commands interact only with the virtual filesystem. Content is loaded from
`CONTENT_PATH`, which defaults to `/app/content`. Vaultsh exposes no commands
that create, modify, or delete files.

## Quickstart

Run Vaultsh with Atlas and their shared content through the sibling `lab`
repository:

```sh
cd ../lab
docker compose up --build
```

Open http://localhost:8080/testui/.

For local development, install Go 1.24 or newer:

```sh
CONTENT_PATH=../lab/content go run ./cmd/vaultsh
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

Portfolio content is stored in the sibling `lab` repository and mounted
read-only at runtime. Vaultsh remains only the virtual filesystem and shell
engine.

Each non-empty line uses a lowercase `key: value` format. Keys may repeat and
blank lines may separate sections.

```text
focus: backend services
focus: distributed systems

technology: Go
technology: Docker
```

## Project Structure

```text
cmd/vaultsh/          application entry point
internal/command/     shell commands
internal/filesystem/  read-only virtual filesystem
internal/httpapi/     HTTP transport
internal/parser/      tokenizer, lexer, and parser
internal/shell/       execution engine and sessions
internal/storage/     mounted-content loader
testui/               browser terminal
```

Shared content and local orchestration documentation live in the `lab`
repository.
