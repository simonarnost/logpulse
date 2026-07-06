# LogPulse

A distributed log aggregator + alerting system, built to learn Go, gRPC, and Protobuf.

Agents tail logs and stream them to a central Collector, which stores entries and fans them out to subscribers. An Alerter watches the stream and fires webhooks on rules like *"5+ ERRORs in 60s from service X"*.

## Status

Early stage — learning project.

| Component | What it does | State |
|-----------|--------------|-------|
| **Agent** | Streams log entries to the Collector (client streaming) | 🚧 mocked sender |
| **Collector** | Receives from N agents, counts/logs entries | ✅ ingest works |
| **Alerter** | Evaluates rules, fires webhooks | ⏳ planned |
| **CLI** (`logpulse`) | `tail`, `query`, `rules ...` | ⏳ planned |

## gRPC patterns

- Client streaming (agent → collector) — implemented
- Server streaming (live tail, alert feed) — planned
- Bidi streaming (heartbeats) — planned
- Unary RPCs (query, rule CRUD) — planned

## Layout

```
proto/     # .proto definitions (data model lives here)
gen/       # generated protobuf/gRPC code — never edit
cmd/       # one entrypoint per service (agent, collector, ...)
internal/  # service logic
```

## Setup

```bash
brew install go protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Run

```bash
# regenerate proto code after editing proto/*.proto
protoc --go_out=gen --go-grpc_out=gen proto/*.proto

# terminal 1 — collector (listens on :50051)
go run cmd/collector/main.go

# terminal 2 — agent (currently sends 3 mock entries)
go run cmd/agent/main.go
```

## Roadmap

- Real log tailing in the agent (files + stdin)
- Persistent storage (in-memory → SQLite/BadgerDB)
- Alerter with rule evaluation + webhooks
- `logpulse` CLI
- Stretch: TUI dashboard, retention policies, Prometheus metrics