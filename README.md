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

## Command Reference

| Command | Usage | Purpose |
| --- | --- | --- |
| `help` | `help [command]` | List commands or show command-specific usage |
| `about` | `about` | Describe Vaultsh |
| `pwd` | `pwd` | Print the current virtual directory |
| `ls` | `ls [-alR] [path]` | List files and directories |
| `cd` | `cd [directory]` | Change the current virtual directory |
| `cat` | `cat [-n] [file]` | Print a file or pipeline input |
| `tree` | `tree [-L depth] [path]` | Print a directory tree |
| `grep` | `grep [-in] <pattern> [file]` | Filter lines using a regular expression |
| `head` | `head [-n count] [file]` | Print the first lines |
| `tail` | `tail [-n count] [file]` | Print the last lines |
| `wc` | `wc [-lwc] [file]` | Count lines, words, and bytes |
| `sort` | `sort [-r] [file]` | Sort lines |
| `history` | `history` | List commands from the current session |
| `clear` | `clear` | Clear the terminal through a backend action |

Directories end with `/` in standard `ls` output. `ls -a` includes hidden
entries, while `ls -l` includes the read-only mode and file size.

List the root with hidden entries and long formatting:

```sh
ls -la /
```

Example output:

```text
-r--r--r--      428 about.txt
-r--r--r--      120 education.txt
dr-xr-xr-x        - experience/
dr-xr-xr-x        - projects/
-r--r--r--      537 skills.txt
```

File sizes change when embedded content is edited.

List recursively:

```sh
ls -R experience
```

Limit a tree to two levels:

```sh
tree -L 2 /
```

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

Print numbered content:

```sh
cat -n education.txt
```

Example output:

```text
     1	institution: University of Rijeka
     2	degree: Master of Education (MEd)
     3	field: Information Technology
     4	graduation_year: 2019
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

List programming languages alphabetically:

```sh
cat skills.txt | grep "^language:" | sort
```

Show the first three ReversingLabs highlights:

```sh
cat experience/reversinglabs.txt | grep "^highlight:" | head -n 3
```

Count A1 stack groups:

```sh
cat experience/a1.txt | grep "^stack:" | wc -l
```

Show numbered, case-insensitive matches:

```sh
cat skills.txt | grep -in "python"
```

Display the five most recent history entries:

```sh
history | tail -n 5
```

Reverse-sort the backend skills:

```sh
cat skills.txt | grep "^backend:" | sort -r
```

Commands that accept `[file]` can read either a virtual file directly or
pipeline input:

```sh
head -n 5 skills.txt
cat skills.txt | head -n 5
```

`grep` returns exit code `1` when no lines match. Syntax and option errors return
exit code `2`. Unknown commands return `127`.

## Keyboard Shortcuts

- `Tab`: complete commands and virtual filesystem paths
- `Ctrl+L`: run the backend-owned `clear` command

Autocomplete uses the current session directory:

```text
cat exp<Tab>       -> cat experience/
cd experience
cat rev<Tab>       -> cat reversinglabs.txt
```

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

Continue in the same directory:

```sh
curl -X POST http://localhost:8080/api/exec \
  -H "Content-Type: application/json" \
  -d '{"line":"pwd","session_id":"<session-id>"}'
```

Expected command output:

```text
/experience
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
ls -R experience
cat skills.txt | sort
cat skills.txt | sort -r
cat skills.txt | head -n 5
cat skills.txt | tail -n 5
cat skills.txt | wc -l
```

The following command options are planned but not implemented yet:

```sh
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
- [x] `ls -R [path]`
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
