# Command Reference

Append `--verbose` to a command or pipeline to include execution metadata in
the API result. The command's normal text output is unchanged.

## Commands

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

Use `help <command>` inside Vaultsh for command-specific usage.

## Filesystem

Directories end with `/` in standard `ls` output. `ls -a` includes hidden
entries, `ls -l` includes read-only mode and byte size, and `ls -R` recursively
lists directories.

```sh
pwd
ls
ls -la /
ls -R experience
tree
tree -L 2 /
cd experience
cat reversinglabs.txt
cd ..
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

Number file or pipeline lines:

```sh
cat -n education.txt
cat education.txt | cat -n
```

## Pipelines

Each pipeline stage receives the previous stage's output. Execution stops when
a stage returns a non-zero exit code.

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

Reverse-sort backend skills:

```sh
cat skills.txt | grep "^backend:" | sort -r
```

Commands that accept `[file]` can read a virtual file directly or consume
pipeline input:

```sh
head -n 5 skills.txt
cat skills.txt | head -n 5
```

More pipeline examples:

```sh
cat skills.txt | grep "^language:" | sort | head -n 5
cat experience/reversinglabs.txt | grep "^stack:" | wc -l
tree | grep ".txt" | sort
```

## Keyboard Shortcuts

- `Tab`: complete commands and virtual filesystem paths
- `Ctrl+L`: run the backend-owned `clear` command

Autocomplete uses the current session directory:

```text
cat exp<Tab>       -> cat experience/
cd experience
cat rev<Tab>       -> cat reversinglabs.txt
```

## Exit Codes

- `0`: success
- `1`: command failure or no `grep` matches
- `2`: syntax, option, or usage error
- `126`: recognized but unsupported behavior
- `127`: command not found

## Planned Options

```sh
ls -lt
```
