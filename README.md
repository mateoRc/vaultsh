# vaultsh
A shell engine featuring a virtual filesystem, command parser, and extensible command system.

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

---

# v2

## Parser
- [x] Tokenizer
- [x] Lexer
- [x] Parser
- [ ] AST

## Pipes
- [ ] |
- [ ] Pipeline executor

## Built-ins
- [ ] grep
- [ ] head
- [ ] tail
- [ ] wc
- [ ] sort

## Command Options
- [ ] `ls -R [path]`
- [ ] `tree -L <depth> [path]`
- [ ] `cat -n <file>`
- [ ] `grep -i` and `grep -n`
- [ ] `head -n <count>` and `tail -n <count>`
- [ ] `sort -r`
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

- [ ] Autocomplete
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
