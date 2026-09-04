# geth-indexer

[![Go Reference](https://pkg.go.dev/badge/github.com/gNimit/geth-indexer.svg)](https://pkg.go.dev/github.com/gNimit/geth-indexer)
[![Go Report Card](https://goreportcard.com/badge/github.com/gNimit/geth-indexer)](https://goreportcard.com/report/github.com/gNimit/geth-indexer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**geth-indexer** is a high-throughput, enterprise-grade Go library and daemon for indexing EVM smart contract events. It serves as an abstraction layer over raw Ethereum JSON-RPC nodes, providing multi-provider RPC pooling, automatic rate limiting, circuit breaking, recursive range bisection, atomic checkpointing, and pluggable sink dispatching.

Designed to exceed **2,000+ events/second** (benchmarked at **1.6M+ events/second** in-memory pipeline capacity), `geth-indexer` is built for mission-critical Web3 data infrastructure.

---

## Key Features

- 🚀 **High Throughput:** Concurrent worker pipeline with zero-allocation ABI decoding and adaptive chunk sizing.
- 🔌 **Importable Go Package:** Simple, functional-option API that returns a native Go channel (`<-chan *types.Event`) for streaming.
- 🔄 **Multi-Provider RPC Pool:** Intelligent load balancing (Least In-Flight, Latency Weighted, Round Robin) across Infura, Alchemy, QuickNode, or local Geth/Reth nodes.
- 🛡️ **Rate Limit & Circuit Breakers:** Per-endpoint token bucket rate limiting and automatic 3-state circuit breakers to handle HTTP 429 ("Too Many Requests") and node failures without downtime.
- 🧩 **Pluggable Event Sinks:** Route events to Go channels, PostgreSQL (idempotent multi-row batch inserts), HTTP Webhooks, or custom sinks (Kafka, gRPC, Redis).
- 💾 **State Tracking & Gap Recovery:** Tracks progress and recovers missed historical block ranges via atomic JSON files or SQL databases.
- 📊 **Prometheus & Grafana Observability:** Built-in Prometheus metrics exporter for throughput, RPC latency, current block, and sink performance.

---

## Installation

```bash
go get github.com/gNimit/geth-indexer
```

---

## Quickstart: Library Usage

### 1. Stream Contract Events via Go Channel

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/indexer"
	"github.com/gNimit/geth-indexer/pkg/rpc"
)

func main() {
	ctx := context.Background()

	// 1. Setup multi-provider RPC pool
	pool, err := rpc.NewPool(rpc.PoolConfig{
		ChainID: 1,
		Endpoints: []rpc.EndpointConfig{
			{URL: "https://eth.llamarpc.com", Provider: "llama", RPS: 25},
			{URL: "https://rpc.ankr.com/eth", Provider: "ankr", RPS: 25},
		},
		Balancer: rpc.NewLeastInFlightBalancer(),
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Setup persistent checkpoint store
	store, _ := checkpoint.NewFileStore("checkpoint.json")

	// 3. Configure indexer
	cfg := indexer.Config{
		ChainID:        1,
		StartBlock:     19000000,
		EndBlock:       0, // 0 = continuous live tailing
		BlockChunkSize: 1000,
		Concurrency:    8,
		Contracts: []indexer.ContractConfig{
			{
				Address: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), // USDC
				ABIJSON: `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`,
				Events:  []string{"Transfer"},
			},
		},
	}

	idx, err := indexer.New(cfg,
		indexer.WithRPCPool(pool),
		indexer.WithCheckpointStore(store),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Start streaming events through channel
	eventsCh, err := idx.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for event := range eventsCh {
		fmt.Printf("[%s] Block %d | Tx %s | Data: %v\n",
			event.EventName, event.BlockNumber, event.TxHash.Hex(), event.Data)
	}
}
```

---

### 2. Index Directly to PostgreSQL

```go
pgSink, err := sink.NewPostgresSink(sink.PostgresConfig{
    Host:      "localhost",
    Port:      5432,
    User:      "postgres",
    Password:  "postgres",
    Database:  "geth_indexer",
    TableName: "indexed_events",
})

idx, err := indexer.New(cfg,
    indexer.WithRPCPool(pool),
    indexer.WithCheckpointStore(store),
    indexer.WithSink(pgSink),
)
```

---

## Quickstart: CLI & Docker

### Running via CLI
```bash
# Build binary
make build

# Run indexer daemon
./bin/geth-indexer \
  -rpc="https://eth.llamarpc.com,https://rpc.ankr.com/eth" \
  -chain-id=1 \
  -contract="0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" \
  -abi="path/to/contract.abi" \
  -from=19000000 \
  -to=0 \
  -concurrency=8 \
  -metrics-addr=":9090"
```

### Running with Docker Compose (PostgreSQL + Prometheus)
```bash
docker-compose up -d
```

Scrape metrics at `http://localhost:9090/metrics` or browse Prometheus at `http://localhost:9091`.

---

## Benchmarks

Benchmarked on Apple M2 (arm64, macOS):

| Benchmark Test | Operations | Latency | Throughput |
|---|---|---|---|
| **ABI Event Decoding** | 2,971,256 ops | 393 ns/op | **2,543,800 events/sec** |
| **Full Pipeline Throughput** | 5,976,000 events | 1.23 ms / 2K events | **1,614,000 events/sec** |

---

## Architecture & Maintainers Guide

- Read [Architecture.md](file:///Users/adrem/Dev/geth-indexer/Architecture.md) for detailed pipeline internals, RPC load balancing, and failure modes.
- Read [Agents.md](file:///Users/adrem/Dev/geth-indexer/Agents.md) for AI agents and maintainer coding invariants and extension guides.

---

## License

MIT License. See [LICENSE](file:///Users/adrem/Dev/geth-indexer/LICENSE) for details.
