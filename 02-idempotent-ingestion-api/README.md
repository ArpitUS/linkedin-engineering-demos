# Idempotent Ingestion API

Run `go test ./... && go run .`.

The in-memory map stands in for a database table keyed by `Idempotency-Key`. The concurrent test simulates duplicate requests arriving at the same time, a common result of a client retrying through a timeout or gateway failure.

Production implementation notes:

- Make the idempotency key unique in the same transactional store as the side effect.
- Persist the original response so repeat requests receive a consistent result.
- Bind the key to the authenticated tenant and request fingerprint.
- Expire keys only after the supported retry window.

Practical takeaway: use a database uniqueness guarantee, not a process-local map, to make duplicate suppression survive replicas and restarts.
