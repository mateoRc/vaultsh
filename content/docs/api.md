# HTTP API

Vaultsh exposes JSON endpoints for command execution and autocomplete.

## Health

```http
GET /healthz
```

Successful response:

```text
ok
```

## Execute a Command

```http
POST /api/exec
Content-Type: application/json
```

Request:

```json
{
  "line": "pwd",
  "session_id": ""
}
```

Response:

```json
{
  "output": "/",
  "exit_code": 0,
  "session_id": "<session-id>"
}
```

The first request can omit `session_id`. The server generates one and returns
it. Send that value with later requests to preserve the working directory and
command history.

```sh
curl -X POST http://localhost:8080/api/exec \
  -H "Content-Type: application/json" \
  -d '{"line":"cd experience"}'
```

Continue the session:

```sh
curl -X POST http://localhost:8080/api/exec \
  -H "Content-Type: application/json" \
  -d '{"line":"pwd","session_id":"<session-id>"}'
```

Sessions expire after 30 minutes of inactivity.

Some commands return a frontend action:

```json
{
  "output": "",
  "exit_code": 0,
  "action": "clear",
  "session_id": "<session-id>"
}
```

Add `--verbose` as the final argument to return execution metadata without
changing command output:

```json
{
  "line": "cat skills.txt | grep Go --verbose"
}
```

```json
{
  "output": "language: Go",
  "exit_code": 0,
  "verbose": "pipeline=cat,grep; stages=2; completed=2",
  "session_id": "<session-id>"
}
```

The `verbose` field is omitted from normal responses. Commands do not receive
the global flag.

## Complete Input

```http
POST /api/complete
Content-Type: application/json
```

Request:

```json
{
  "line": "cat exp",
  "cursor": 7,
  "session_id": "<session-id>"
}
```

Response:

```json
{
  "start": 4,
  "end": 7,
  "replacement": "experience/",
  "candidates": [
    "experience/"
  ],
  "session_id": "<session-id>"
}
```

`start` and `end` identify the input range the client should replace.
`replacement` is the longest common candidate prefix.

## Errors

- Invalid JSON returns HTTP `400`.
- Session creation failure returns HTTP `500`.
- Shell command failures are returned as HTTP `200` with a non-zero
  `exit_code`.
