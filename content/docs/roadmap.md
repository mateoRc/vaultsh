# Roadmap

## MVP

### Bootstrap

- [x] Containerized Go HTTP server
- [x] Docker Compose
- [x] Disposable test UI
- [x] Health endpoint

### HTTP API

- [x] `POST /api/exec`
- [x] Request/response models
- [x] Unknown command handling

### Shell Engine

- [x] Engine
- [x] Command registry

### Commands

- [x] `help`
- [x] `about`
- [x] `clear` backend action
- [x] Command interface

### Frontend

- [x] Send commands to API
- [x] Render output
- [x] Basic history
- [x] Auto focus
- [x] Clear shortcut (`Ctrl+L`)

### Developer Experience

- [x] Graceful shutdown
- [x] Structured logging
- [x] Basic tests
- [x] Embedded content format validation
- [x] GitHub Actions CI

## v1

### Virtual Filesystem

- [x] Node model
- [x] Directory
- [x] File
- [x] Path resolver
- [x] Current working directory

### Commands

- [x] `pwd`
- [x] `ls`
- [x] `cd`
- [x] `cat`
- [x] `tree`
- [x] `history`
- [x] Command-specific help

### Sessions

- [x] Execution context
- [x] Session ID
- [x] Command history
- [x] Working directory per session
- [x] Session expiration and cleanup

## v2

### Parser

- [x] Tokenizer
- [x] Lexer
- [x] Parser
- [x] AST

### Pipelines

- [x] Pipe operator
- [x] Pipeline executor

### Built-ins

- [x] `grep`
- [x] `head`
- [x] `tail`
- [x] `wc`
- [x] `sort`

### Command Options

- [x] `ls -a`, `ls -l`, and combined forms
- [ ] File timestamps and `ls -t`
- [x] `ls -R [path]`
- [x] `tree -L <depth> [path]`
- [x] `cat -n <file>`
- [x] `grep -i` and `grep -n`
- [x] `head -n <count>` and `tail -n <count>`
- [x] `sort -r`
- [x] Verbose output (`--verbose`) returned with command results

## v3

### Storage

- [x] Embedded text content
- [ ] Session store abstraction
- [ ] External session store
- [ ] Filesystem abstraction

### Authentication

- [ ] JWT authentication

### Performance

- [ ] LRU cache
- [ ] Benchmarks
- [ ] Profiling

## v4

### Search

- [ ] Search service
- [ ] Index
- [ ] Ranking
- [ ] Query language

## Nice to Have

- [x] Command and path autocomplete
- [ ] Aliases
- [ ] Environment variables
- [ ] Configuration
- [ ] Plugin system
- [ ] Multiple mounted vaults
- [ ] WebSocket transport
- [ ] TUI client
- [ ] CLI client
- [ ] Fuzzy history search with `Ctrl+R`

### Easter Eggs

- [x] `.motd` hidden welcome message
- [x] `whoami` prints Mateo's current role and focus
- [ ] Single-command triggers
- [ ] Command-sequence triggers
- [ ] Session-scoped progress
- [x] `hire mateo` returns a sudo-access hint
- [x] `sudo hire mateo -s <salary>` requires a salary
- [x] Return a counter-offer at 1.5× the submitted salary
- [x] Accept a pending counter-offer with `y` or `Y`
