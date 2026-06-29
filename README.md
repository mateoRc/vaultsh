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
- [ ] Directory
- [ ] File
- [ ] Path resolver
- [ ] Current working directory

## Commands
- [ ] pwd
- [ ] ls
- [ ] cd
- [ ] cat
- [ ] tree

## Sessions
- [ ] Execution context
- [ ] Session ID
- [ ] Command history
- [ ] Working directory per session

---

# v2

## Parser
- [ ] Tokenizer
- [ ] Lexer
- [ ] Parser
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
- [ ] Verbose output (`--verbose`) returned with command results

---

# v3

## Storage
- [ ] Embedded content
- [ ] SQLite backend
- [ ] Filesystem abstraction

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
