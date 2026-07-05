FROM golang:1.24-alpine AS test

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go test ./...

FROM test AS build

RUN CGO_ENABLED=0 go build -o /vaultsh ./cmd/vaultsh

FROM alpine:3.21

ARG VERSION=dev
LABEL org.opencontainers.image.title="Vaultsh" \
      org.opencontainers.image.version="${VERSION}"

WORKDIR /app

COPY --from=build /vaultsh /usr/local/bin/vaultsh
COPY web ./web

EXPOSE 8080

CMD ["vaultsh"]
