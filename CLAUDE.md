# LogPulse

Distributed log aggregator + alerting system. Learning project for Go, gRPC, and Protobuf.

## Architecture

3 services + 1 CLI:

- **Agent** — tails log files / reads stdin, streams entries to Collector (client streaming)
- **Collector** — receives from N agents, stores logs (in-memory or SQLite), fans out to subscribers (bidi streaming w/ heartbeats)
- **Alerter** — subscribes to Collector stream, evaluates rules (e.g. "5+ ERRORs in 60s from service X"), fires webhooks (server streaming for alert notifications)
- **CLI** (`logpulse`) — `tail`, `query`, `rules add/list/delete`

## gRPC patterns used

- Client streaming (agent → collector)
- Server streaming (live tail, alert feed)
- Bidi streaming (heartbeats)
- Unary RPCs (query, CRUD alert rules)
- Interceptors (request logging, auth)
- Proto features: `oneof`, enums, `google.protobuf.Timestamp`, nested messages

## Project structure

```
logpulse/
├── proto/              # .proto definitions (start here)
├── cmd/
│   ├── agent/          # agent entrypoint
│   ├── collector/      # collector entrypoint
│   ├── alerter/        # alerter entrypoint
│   └── logpulse/       # CLI entrypoint
├── internal/
│   ├── agent/
│   ├── collector/
│   ├── alerter/
│   └── common/         # shared types, interceptors
├── pkg/                # public reusable packages (if any)
└── gen/                # generated protobuf/gRPC code
```

## Setup

```bash
# Install Go (1.22+)
brew install go

# Install protobuf toolchain
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Init module
go mod init github.com/simonarnost/logpulse
```

## Build & run

```bash
# Generate proto code
protoc --go_out=gen --go-grpc_out=gen proto/*.proto

# Run services
go run cmd/collector/main.go
go run cmd/alerter/main.go
go run cmd/agent/main.go --file /var/log/system.log

# CLI
go run cmd/logpulse/main.go tail
go run cmd/logpulse/main.go query --level ERROR --since 1h
go run cmd/logpulse/main.go rules add --name "error-spike" --pattern "level=ERROR" --threshold 5 --window 60s
```

## Development

- Start with proto definitions in `proto/` — design the data model first
- Generated code goes in `gen/` — never edit generated files
- Each service is its own `cmd/` entrypoint with logic in `internal/`
- Use `buf` instead of raw `protoc` once comfortable (better DX)

## Stretch goals

- TUI dashboard (bubbletea)
- Log retention policies
- Multi-collector fanout
- Persistent storage (SQLite/BadgerDB)
- Prometheus metrics endpoint
