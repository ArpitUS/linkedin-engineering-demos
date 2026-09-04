# Versioned API Contracts

Run `go test ./... && go run .`.

The V1 response is treated as a permanent contract. V2 adds `displayName` and `tier` without renaming or changing the semantics of `id` and `name`. The V1 test is a contract test: it protects a downstream client from an accidental breaking change.

Before removing a version, publish a deprecation date, capture client/version usage in telemetry, provide a migration guide, and verify that active consumers have moved. For a compatible field addition, prefer evolving the existing version and make clients tolerate unknown fields.

Practical takeaway: put a consumer-visible contract test beside every public API version, then let usage telemetry determine when removal is safe.
