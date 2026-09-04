package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gNimit/geth-indexer/pkg/metrics"
	"github.com/gNimit/geth-indexer/pkg/rpc"
)

func main() {
	ctx := context.Background()

	// Multi-provider RPC Pool with weights, rate limits, and circuit breakers
	poolConfig := rpc.PoolConfig{
		ChainID: 1,
		Endpoints: []rpc.EndpointConfig{
			{
				URL:              "https://eth.llamarpc.com",
				Provider:         "llama",
				Weight:           2,
				RPS:              25, // 25 requests/sec max
				Burst:            50,
				Timeout:          5 * time.Second,
				FailureThreshold: 3,
				CooldownDuration: 15 * time.Second,
			},
			{
				URL:              "https://rpc.ankr.com/eth",
				Provider:         "ankr",
				Weight:           1,
				RPS:              20,
				Burst:            30,
				Timeout:          5 * time.Second,
				FailureThreshold: 3,
				CooldownDuration: 15 * time.Second,
			},
			{
				URL:              "https://cloudflare-eth.com",
				Provider:         "cloudflare",
				Weight:           1,
				RPS:              15,
				Burst:            20,
				Timeout:          5 * time.Second,
				FailureThreshold: 3,
				CooldownDuration: 15 * time.Second,
			},
		},
		Balancer:       rpc.NewLatencyWeightedBalancer(),
		Metrics:        metrics.NewPrometheusCollector("multi_rpc"),
		HealthInterval: 10 * time.Second,
	}

	pool, err := rpc.NewPool(poolConfig, nil)
	if err != nil {
		log.Fatalf("Failed to create RPC pool: %v", err)
	}
	defer pool.Close()

	// Query latest block height
	blockNum, err := pool.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("Failed to query block number: %v", err)
	}
	fmt.Printf("Latest Ethereum Block: %d\n", blockNum)

	// Fetch logs across pool with automatic rotation
	logs, err := pool.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: common.Big0,
		ToBlock:   common.Big1,
	})
	if err != nil {
		fmt.Printf("Query error (handled): %v\n", err)
	} else {
		fmt.Printf("Fetched %d logs\n", len(logs))
	}

	// Print endpoint health & latency metrics
	for _, ep := range pool.Endpoints() {
		fmt.Printf("Provider: %-12s | Available: %-5v | In-Flight: %d | Latency EMA: %.2fms\n",
			ep.Config().Provider, ep.IsAvailable(), ep.InFlight(), ep.LatencyEMA())
	}
}
