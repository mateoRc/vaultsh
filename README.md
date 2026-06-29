# vaultsh
A shell engine featuring a virtual filesystem, command parser, and extensible command system.

# Startup

## Docker

Prerequisites:
- Docker Engine or Docker Desktop
- Docker Compose v2

Start Vaultsh:

```sh
docker compose up --build
```

Open:

```text
http://localhost:8080/testui/
```

Stop Vaultsh:

```sh
docker compose down
```

## Local Development

Prerequisites:
- Go 1.24 or newer

Run from the repository root:

```sh
go run ./cmd/vaultsh
```

Run tests:

```sh
go test ./...
```

The test UI is plain HTML, CSS, and JavaScript; it does not require Node.js or npm.

# Usage

## Explore the Virtual Filesystem

```sh
help
help cat
pwd
tree
ls
ls -la
ls experience
cd experience
pwd
cat reversinglabs.txt
cd ..
cat skills.txt
history
```

Directories end with `/` in standard `ls` output. `ls -a` includes hidden
entries, while `ls -l` includes the read-only mode and file size.

Paths can be absolute or relative:

```sh
cat /about.txt
cat projects/vaultsh.txt
cd /experience
cd ..
```

Quoted and escaped arguments are supported:

```sh
cat "file with spaces.txt"
cat file\ with\ spaces.txt
```

## Pipelines

Pipeline execution and regular-expression filtering are available:

```sh
cat skills.txt | grep "^language:"
cat skills.txt | grep -i python
cat experience/reversinglabs.txt | grep -n highlight
```

Each stage receives the previous stage's output. Execution stops when a stage
returns a non-zero exit code.

## Keyboard Shortcuts

- `Tab`: complete commands and virtual filesystem paths
- `Ctrl+L`: run the backend-owned `clear` command

## HTTP API

Execute a command:

```sh
curl -X POST http://localhost:8080/api/exec \
  -H "Content-Type: application/json" \
  -d '{"line":"tree"}'
```

The response includes a `session_id`. Send it with later requests to preserve
the working directory and history:

```sh
curl -X POST http://localhost:8080/api/exec \
  -H "Content-Type: application/json" \
  -d '{"line":"cd experience","session_id":"<session-id>"}'
```

## Advanced Examples

Compose multiple pipeline stages:

```sh
cat skills.txt | grep "^language:" | sort
cat skills.txt | grep "^backend:" | sort -r
cat experience/reversinglabs.txt | grep highlight | head -n 3
cat experience/reversinglabs.txt | grep "^stack:" | wc -l
history | tail -n 5
tree | grep ".txt" | sort
```

Limit and inspect output:

```sh
cat -n about.txt
tree -L 2
cat skills.txt | sort
cat skills.txt | sort -r
cat skills.txt | head -n 5
cat skills.txt | tail -n 5
cat skills.txt | wc -l
```

The following command options are planned but not implemented yet:

```sh
ls -R experience
ls -lt
```

# Roadmap

## MVP

### Bootstrap
- [x] Containerized Go HTTP server
- [x] Docker Compose
- [x] Disposable test UI
- [x] Health endpoint

### HTTP API
- [x] POST /api/exec
- [x] Request/response models
- [x] Unknown command handling

### Shell Engine
- [x] Engine
- [x] Command registry

### Commands
- [x] help
- [x] about
- [x] clear (backend action)
- [x] Promote commands to Command interface

### Frontend (Test UI)
- [x] Send commands to API
- [x] Render output
- [x] Basic history
- [x] Auto focus
- [x] Clear shortcut (Ctrl+L)

### Developer Experience
- [x] Graceful shutdown
- [x] Structured logging
- [x] Basic tests

---

# v1

## Virtual Filesystem
- [x] Node model
- [x] Directory
- [x] File
- [x] Path resolver
- [x] Current working directory

## Commands
- [x] pwd
- [x] ls
- [x] cd
- [x] cat
- [x] tree
- [x] history
- [x] Command-specific help

## Sessions
- [x] Execution context
- [x] Session ID
- [x] Command history
- [x] Working directory per session
- [x] Session expiration and cleanup

---

# v2

## Parser
- [x] Tokenizer
- [x] Lexer
- [x] Parser
- [x] AST

## Pipes
- [x] |
- [x] Pipeline executor

## Built-ins
- [x] grep
- [x] head
- [x] tail
- [x] wc
- [x] sort

## Command Options
- [x] `ls -a`, `ls -l`, and combined forms
- [ ] File timestamps and `ls -t`
- [ ] `ls -R [path]`
- [x] `tree -L <depth> [path]`
- [x] `cat -n <file>`
- [x] `grep -i` and `grep -n`
- [x] `head -n <count>` and `tail -n <count>`
- [x] `sort -r`
- [ ] Verbose output (`--verbose`) returned with command results

---

# v3

## Storage
- [x] Embedded content
- [ ] Session store abstraction
- [ ] External session store
- [ ] SQLite backend
- [ ] Filesystem abstraction

## Authentication
- [ ] JWT authentication

## Performance
- [ ] LRU cache
- [ ] Benchmarks
- [ ] Profiling

---

# v4

## Search
- [ ] Search service
- [ ] Index
- [ ] Ranking
- [ ] Query language

---

# Nice to Have

- [x] Command and path autocomplete
- [ ] Aliases
- [ ] Environment variables
- [ ] Configuration
- [ ] Plugin system
- [ ] Multiple mounted vaults
- [ ] WebSocket transport
- [ ] TUI client
- [ ] CLI client
- [ ] Fuzzy history search with Ctrl+R (frontend)

## Easter Eggs
- [ ] Single-command triggers
- [ ] Command-sequence triggers
- [ ] Session-scoped progress (for example, 5 out of 5)
