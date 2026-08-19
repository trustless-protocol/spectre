# Metrics and best practices

The goal of metric collection is to provide data-driven insights into the inner workings of your services. 

The particular goals we aim to achieve, and that you should keep in mind while implementing metrics, are:
- Debugging
  - Rapid issue identification through structured error metrics and performance bottlenecks
- Operational Health
  - Proactive monitoring of service availability and resource utilization  
- Performance Optimization
  - Data-driven insights into latency patterns and throughput characteristics
- Alerting
  - Automated incident response based on service-level indicators

## Key principles

- Every metric has a purpose
  - Every metric needs to be documented as to why it exists, and what it is meant to capture (business or technical goal).
- RED metrics
  - Start with RED metrics (Rate, Errors, Duration) for every service
- Correlate with logs and traces
  - Correlate metrics with logs and traces using correlation IDs for effective debugging
- Alert on symptoms, not causes
  - Focus on metrics that have an impact, not internal system states that may or may not cause any real-world impact
- Open Standards
  - Use open standards like OpenTelemetry for vendor neutrality
- Include metrics in the design phase
  - Implement observability by design - include metrics requirements from service conception, don't just tack 'em on later.
- Use established conventions and standards
  - Use consistent naming conventions across all services _and_ teams

## What to measure

- Request lifecycle
  - Every request, message, or transaction moving through your system. 
  - Track the full lifecycle: received, processing, succeeded, failed, timed out.
- Connection health
  - Track whether your service can talk to what it needs to: internal services, 3rd party services, RPC nodes, databases. 
  - Include connection state, latency to dependencies, and time until critical resources expire (like IBC clients).
- Error patterns
  - Group errors into categories that make sense for debugging: timeout errors, validation failures, insufficient funds, and network issues. 
  - You want to spot patterns, not track every unique error message.
- Performance bottlenecks
  - Request latencies at different stages of processing, queue depths, batch sizes, and throughput rates. 
  - Focus on the operations that directly impact user experience or system stability.
- Resource consumption
  - Memory usage, CPU utilization, goroutines/threads, file descriptors.
  - Track both your service's consumption and any limits you're operating under.

## When to measure

### Add metrics when:
- You're implementing a new critical path that could fail
- You've had an incident and realized you were flying blind
- You need to validate a critical performance assumption
- You're adding SLAs or need to track service level objectives
- There's a resource that could be exhausted (connections, memory, etc.)

### Skip metrics when:
- The information is already available through logs or traces
- You're tracking something "just in case" with no clear use case
- The cardinality would be unbounded (user IDs, transaction hashes, timestamps)
- It duplicates an existing metric at a slightly different granularity

## Conventions and best practices

### Metric names
Use this format: `namespace_subsystem_name_unit`

Good examples:
- `solidity_ibc_relayer_packets_sent_total`
- `solidity_ibc_relayer_client_expiry_seconds`
- `solidity_ibc_relayer_rpc_duration_seconds`

Bad examples:
- `PacketsSent` - wrong case, no namespace
- `relayer_memory_mb` - use base units (bytes, not MB)
- `connection_count` - should be `connections_total` if it's a counter

### Labels

> _Cardinality_ is about how many unique values a metric's labels (or dimensions) can take.
> For example: if you have a metric for HTTP requests with the labels `status_code` and `method`,
> and `status_code` has 5 possible values and `method` has 3 possible values, the cardinality for those combined labels is 5 × 3 = 15 unique combinations.

- Labels should have bounded, predictable values:
  - `chain_id`: "cosmoshub-4", "osmosis-1"
  - `status`: "success", "timeout", "invalid_proof"
- Keep cardinality under 100 per metric
  - If you can't list all possible values, don't make it a label!

### Other rules

- Never expose sensitive data
  - No private keys, passwords, or internal IPs in metrics
- Always use seconds for time - Not milliseconds, not minutes. Seconds.
- Counters only go up
  - Never reset a counter to zero. Let Prometheus handle resets.
- One metric, one meaning
  - Don't reuse the same metric name for different things
- Base units only
  - Bytes not kilobytes, seconds not milliseconds

## Metric types

### Counter
A number that only goes up. Perfect for counting things that happen - requests, errors, bytes processed. You'll usually query these with `rate()` or `increase()`.

When to use: Total requests served, packets sent, errors encountered, and bytes transferred.

### Gauge
A value that can go up or down. Current state of something - active connections, queue depth, temperature, and available memory.

When to use: Active connections, queue sizes, current memory usage, and last successful run timestamp.

### Histogram
Sample observations and count them in buckets. Great for latencies and sizes where you care about percentiles. More expensive than counters and gauges.

When to use: Request durations, response sizes, processing times - anything where you need percentiles.

### Summary
Similar to a histogram, but calculates quantiles on the client side. Generally, histograms are preferred because they can be aggregated.

When to use: Rarely. You can use histograms instead unless you have a specific reason.

## Language-specific guides

### Rust

TODO: IBC-143 Write the Rust-specific guide

## Relayer queue-state logs

The Go relayer has no metrics exporter yet. Until RLY-02 provides one, every
services engine emits a structured `[QueueState]` line at startup and once per
minute. These values are deliberately bounded gauges so a log pipeline can
alert on them today and map them directly to an exporter later.

```
[QueueState] queued(cosmos=0 eth=0) pending(cosmos=2 eth=0 l2=0) timeout-dead-letter(cosmos=0 eth=0 l2=0)
```

| Field | Future metric type | Meaning |
|---|---|---|
| `queued(cosmos,eth)` | Gauge | Packets still waiting for the batch adapters to flush. Sustained growth means work is arriving faster than delivery. |
| `pending(cosmos,eth,l2)` | Gauge | Sent packets still eligible for timeout recovery. |
| `timeout-dead-letter(cosmos,eth,l2)` | Gauge | Timeout submissions that exhausted the bounded retry budget. A non-zero value means escrowed funds need operator action. |

Two `[QueueState][ATTENTION]` lines are emitted independently of the normal
state line. The first reports timeout dead letters and their oldest age. The
second reports packets repeatedly deferred by a shared prerequisite, including
the worst deferral count. Deferrals never consume the timeout retry budget, so
they do not dead-letter; alerting on that line is required to distinguish a
healthy empty queue from an RPC or light-client prerequisite that has stalled.

The reporter stops with its adapter context. Pending timeout entries and dead
letters are persisted as atomically replaced JSON snapshots under
`.fast-ibc-state/` beside the selected relayer config and are restored before
startup recovery begins. A corrupt or unwritable state location aborts startup
instead of silently granting packets a fresh retry budget. A timeout dead
letter remains visible and suppressed for 30 days, then automatically returns
to timeout recovery with a fresh retry budget; this prevents a stale
client/sequence-number reuse from being shadowed forever. RLY-02 is still
responsible for an operator-facing drain/reactivation workflow.
