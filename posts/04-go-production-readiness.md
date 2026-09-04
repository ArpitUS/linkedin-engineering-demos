# What I review before approving a Go service for production

![Go production workflow](graphics/04-go-production-readiness.svg)

For a Go service, “the tests pass” is not a production-readiness signal. I want to know what happens when dependencies stall, traffic bursts, and a deployment asks the process to stop.

Three checks catch a large class of operational failures:

**Deadline propagation**

Every dependency call needs the request context and a service-owned deadline. Otherwise a slow downstream system can retain goroutines, sockets, and memory indefinitely.

**Bounded concurrency**

Unlimited goroutines are not backpressure. A bounded worker pool fails fast when the service is saturated, allowing callers, load balancers, and autoscaling systems to see the pressure rather than hiding it in an ever-growing queue.

**Predictable shutdown**

Once shutdown begins, stop accepting new work. Let in-flight work finish only within its deadline. Readiness should become false before termination so Kubernetes can remove the pod from traffic.

The runtime test suite simulates a dependency that never returns. It verifies that the dependency receives cancellation, a second request gets `ErrBusy` rather than queueing forever, and shutdown rejects new work.

Run the demonstration:

```text
cd 04-go-production-readiness
```

Source: https://github.com/ArpitUS/linkedin-engineering-demos/tree/master/04-go-production-readiness

Practical takeaway: every goroutine, dependency call, and work queue must have a limit, a cancellation path, and a metric that exposes when the limit is reached.
