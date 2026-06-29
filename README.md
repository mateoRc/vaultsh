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
- [ ] Engine
- [ ] Execution context
- [ ] Command registry
- [ ] Command interface

### Commands
- [ ] help
- [ ] about
- [ ] clear (frontend only)

### Frontend (Test UI)
- [ ] Send commands to API
- [ ] Render output
- [ ] Basic history
- [ ] Auto focus

### Developer Experience
- [ ] Graceful shutdown
- [ ] Structured logging
- [ ] Basic tests

---

# v1

## Virtual Filesystem
- [ ] Node model
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
