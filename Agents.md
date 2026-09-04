# Agents Guide: Geth-Indexer

Welcome, AI Agent / Contributor. This document provides engineering principles, architectural invariants, code conventions, and maintenance protocols for `geth-indexer`.

---

## 1. Core Engineering Principles & Invariants

When modifying or extending `geth-indexer`, you MUST adhere to the following rules:

1. **Throughput Invariant:**
   - Any modifications to the hot path (`internal/abi`, `pkg/indexer/pipeline.go`, `pkg/sink`) must maintain a minimum processing speed of $\ge 2,000$ events/second.
   - Run `go test -v -bench=. -benchmem ./...` before submitting any code changes.
2. **Deadlock Prevention on Channels:**
   - Never perform unbounded blocking operations on native Go channels (`s.ch <- event`) without monitoring both `ctx.Done()` and `stopCh`.
   - Always ensure consumers can terminate without leaving producer goroutines stranded in memory.
3. **No Unsafe SQL Concatenation:**
   - Never construct SQL queries with `fmt.Sprintf` for variable parameters. All SQL operations in `pkg/checkpoint` and `pkg/sink` must utilize parameterized queries (`$1, $2, ...`) to eliminate SQL injection risks.
4. **Idempotent Checkpoints:**
   - Checkpoint updates must be strictly monotonic: a lower block number must never overwrite a higher indexed block number.
   - File-based checkpoint stores must employ atomic file writing (`.tmp` write followed by `os.Rename`).
5. **Thread Safety:**
   - All public APIs (`rpc.Pool`, `indexer.Indexer`, `checkpoint.Store`, `sink.Sink`) must be safe for concurrent access by multiple goroutines. Protect shared state with `sync.RWMutex` or atomic primitives.

---

## 2. Directory Layout & Subsystem Responsibilities

```
├── cmd/
│   └── geth-indexer/          # Production CLI binary entry point
├── pkg/
│   ├── types/                 # Universal data models (Event, BlockRange)
│   ├── indexer/               # Engine pipeline, coordinator, worker pool, options
│   ├── rpc/                   # Multi-provider RPC pool, circuit breaker, balancers
│   ├── checkpoint/            # State stores (FileStore, SQLStore, MemoryStore)
│   ├── sink/                  # Consumer sinks (Channel, Postgres, Webhook, Stdout)
│   └── metrics/               # Prometheus / OpenTelemetry telemetry abstraction
├── internal/
│   └── abi/                   # Dynamic ABI parser, topic matcher, fast decoder
├── examples/                  # Ready-to-run code examples for library consumers
├── Architecture.md            # System architecture and performance specifications
└── Agents.md                  # This agent guidance manual
```

---

## 3. How to Extend the Codebase

### 3.1 Adding a New Sink (e.g., Kafka, gRPC, Redis)
1. Implement the `sink.Sink` interface in `pkg/sink/`:
   ```go
   type Sink interface {
       Name() string
       Send(ctx context.Context, batch []*types.Event) error
       Close() error
   }
   ```
2. For streaming message brokers (Kafka, RabbitMQ), batch sends inside `Send(ctx, batch)` to maximize I/O throughput.
3. Add a functional option in `pkg/indexer/options.go` (e.g., `WithKafkaSink(...)`).
4. Write unit tests in `pkg/sink/<name>_test.go`.

### 3.2 Adding a New Checkpoint Store (e.g., Redis, DynamoDB)
1. Implement the `checkpoint.Store` interface in `pkg/checkpoint/`:
   ```go
   type Store interface {
       GetLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address) (uint64, error)
       SaveLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address, block uint64) error
       GetUnindexedRanges(ctx context.Context, chainID uint64, contract common.Address) ([]BlockRange, error)
       AddUnindexedRange(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error
       MarkRangeCompleted(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error
       Close() error
   }
   ```
2. Ensure `testStoreContract` in `pkg/checkpoint/store_test.go` passes for your store implementation.

### 3.3 Adding a New RPC Balancer
1. Implement the `rpc.Balancer` interface in `pkg/rpc/balancer.go`:
   ```go
   type Balancer interface {
       Select(endpoints []*Endpoint) (*Endpoint, error)
       Name() string
   }
   ```
2. Verify that unavailable endpoints (`!ep.IsAvailable()`) are skipped.
3. Add a unit test verifying selection behavior in `pkg/rpc/pool_test.go`.

---

## 4. Testing & Verification Protocols

Always run the following commands to verify code changes:

```bash
# 1. Compile all packages and binaries
go build ./...

# 2. Run all unit tests with data race detector
go test -v -race ./...

# 3. Execute performance benchmarks
go test -v -bench=. -benchmem ./...

# 4. Build Docker container
make docker-build
```
