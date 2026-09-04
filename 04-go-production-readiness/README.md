# Go Production Readiness

Run `go test ./...`.

This is a focused service-level demonstration of three production requirements:

- A context deadline is passed to the dependency, so a slow dependency cannot hold a request forever.
- A bounded worker channel supplies backpressure instead of accumulating an unbounded goroutine queue.
- Shutdown rejects new work deterministically.

The tests exercise the edge conditions rather than only the happy path. For an HTTP service, pair this pattern with server read/write timeouts, readiness that becomes false before shutdown, structured error logging, and metrics for deadlines and `ErrBusy` responses.

Practical takeaway: every goroutine, dependency call, and queue needs a defined limit and an observable failure path.
