# Retries Are Not Deduplication

Run `go test ./... && go run .`.

The first path models the important edge condition: the server commits an event but the response is lost. Retrying is correct client behavior, yet a naive ingestion service counts the event twice.

The second path uses an event ID generated before transmission. Replaying that ID becomes a no-op. In a real service, store this identity in a durable database with a retention period that matches the maximum replay window.

Practical takeaway: retries provide at-least-once delivery; a stable event ID plus durable deduplication provides correct aggregates.
