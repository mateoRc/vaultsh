FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd

RUN CGO_ENABLED=0 go build -o /vaultsh ./cmd/vaultsh

FROM alpine:3.21

COPY --from=build /vaultsh /usr/local/bin/vaultsh

EXPOSE 8080

CMD ["vaultsh"]
