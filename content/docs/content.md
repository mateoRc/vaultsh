# Content

Vaultsh exposes the files under `content/` through its read-only virtual
filesystem. The files are embedded into the binary at build time.

## Layout

```text
/
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

## Format

- Use UTF-8 plain text with LF line endings.
- Write one lowercase `key: value` pair per line.
- Repeat keys when a field has multiple values.
- Use blank lines to separate sections.
- Do not use Markdown.
- Do not include confidential or sensitive employer or project details.

Example:

```text
company: Example
role: Backend Engineer

focus: backend services
focus: distributed systems

technology: Go
technology: Docker
```

Run `go test ./...` after editing content. Tests enforce the expected embedded
file layout and format.
