# LinkedIn Engineering Demos

Five small, dependency-free Go projects that turn the accompanying LinkedIn posts into runnable demonstrations. Each folder is intentionally independent: run its commands from that folder.

| Demo | Topic | Run |
| --- | --- | --- |
| `01-retries-and-duplicates` | At-least-once delivery and deduplication | `go test ./... && go run .` |
| `02-idempotent-ingestion-api` | Idempotency keys for write APIs | `go test ./... && go run .` |
| `03-kubernetes-bottleneck-signals` | Correlating latency, throttling, and lag | `go test ./... && go run .` |
| `04-go-production-readiness` | Deadlines, shutdown, and bounded work | `go test ./...` |
| `05-versioned-api-contracts` | Additive contract evolution | `go test ./... && go run .` |

All examples use only the Go standard library, so they are easy to run during a screen recording or live walkthrough.

## Post Series

Publication-ready LinkedIn copy and workflow graphics are in [`posts/`](posts/README.md). Each post links to its matching runtime and test command.
