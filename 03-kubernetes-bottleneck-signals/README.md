# Kubernetes Bottleneck Signals

Run `go test ./... && go run .`.

The executable intentionally distinguishes a backlog from a confirmed bottleneck. Lag alone says that work is accumulating; it does not identify whether the cause is insufficient consumers, a partition skew, a slow database, or a CPU-limited pod.

`prometheus-alerts.yaml` shows an alert that waits for three signals together: high p95 request latency, CPU throttling, and consumer lag. Adjust labels, thresholds, and metric names to match the deployed telemetry stack.

Practical takeaway: page on correlated saturation signals, then use the individual metrics to locate the constrained resource.
