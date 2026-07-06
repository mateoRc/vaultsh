# Vaultsh

A Go service that presents mounted Markdown through a read-only virtual shell.
It provides sessions, pipelines, completion, a browser terminal, and optional
Atlas search and Forge telemetry.

Vaultsh never invokes a host shell or exposes the host filesystem.

## Develop

Requires Go 1.24 or newer and content from the sibling `lab` repository.

```sh
CONTENT_PATH=../lab/content go run ./cmd/vaultsh
go test ./...
```

Run the integrated stack instead:

```sh
cd ../lab
docker compose up --build
```

Open http://localhost:8080/vault/.

## Configuration

- `CONTENT_PATH`: mounted Markdown root; default `/app/content`
- `ATLAS_URL` and `ATLAS_AUTH_TOKEN`: optional search integration
- `FORGE_URL` and `FORGE_AUTH_TOKEN`: optional telemetry integration
- `SESSION_LIMIT`: active-session cap; default `5000`
- `TRUST_PROXY_HEADERS`: trust proxy client-IP headers only when Vaultsh is
  reachable exclusively through a trusted reverse proxy
- `DEPLOYMENT_METADATA_PATH`: optional CI deployment metadata
- `SENTINEL_METADATA_PATH`: optional Sentinel assessment metadata

When an integration URL is set, its token is required.

## Project layout

```text
cmd/vaultsh/          application entry point
internal/command/     shell commands
internal/filesystem/  read-only virtual filesystem
internal/httpapi/     HTTP transport
internal/parser/      tokenizer, lexer, and parser
internal/shell/       execution engine and sessions
internal/storage/     mounted-content loader
web/                  browser terminal
```

The sibling `lab` repository owns shared content, production orchestration,
the command/API references, and
[architecture documentation](https://github.com/mateoRc/lab/tree/main/content/docs).
