package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/indexer"
	"github.com/gNimit/geth-indexer/pkg/rpc"
	"github.com/gNimit/geth-indexer/pkg/sink"
)

const uniswapV3PoolABI = `[
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "sender", "type": "address"},
			{"indexed": true, "name": "recipient", "type": "address"},
			{"indexed": false, "name": "amount0", "type": "int256"},
			{"indexed": false, "name": "amount1", "type": "int256"},
			{"indexed": false, "name": "sqrtPriceX96", "type": "uint160"},
			{"indexed": false, "name": "liquidity", "type": "uint128"},
			{"indexed": false, "name": "tick", "type": "int24"}
		],
		"name": "Swap",
		"type": "event"
	}
]`

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 1. Initialize RPC Pool
	pool, err := rpc.NewPool(rpc.PoolConfig{
		ChainID: 1,
		Endpoints: []rpc.EndpointConfig{
			{URL: "https://eth.llamarpc.com", Provider: "llama", RPS: 25},
		},
	}, nil)
	if err != nil {
		log.Fatalf("RPC Pool initialization failed: %v", err)
	}

	// 2. Configure File-based Checkpointing
	store, err := checkpoint.NewFileStore("uniswap_checkpoint.json")
	if err != nil {
		log.Fatalf("Checkpoint store failed: %v", err)
	}

	// 3. Configure Uniswap USDC/ETH Pool Swap events
	poolAddress := common.HexToAddress("0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640")
	cfg := indexer.Config{
		ChainID:           1,
		StartBlock:        18000000,
		EndBlock:          18010000,
		BlockChunkSize:    1000,
		Concurrency:       8,
		Confirmations:     12,
		SinkBatchSize:     500,
		SinkFlushInterval: 100 * time.Millisecond,
		Contracts: []indexer.ContractConfig{
			{
				Address: poolAddress,
				ABIJSON: uniswapV3PoolABI,
				Events:  []string{"Swap"},
			},
		},
	}

	// 4. Configure PostgreSQL Batch Sink
	pgConfig := sink.PostgresConfig{
		Host:      "localhost",
		Port:      5432,
		User:      "postgres",
		Password:  "postgres",
		Database:  "geth_indexer",
		TableName: "uniswap_swaps",
	}

	idx, err := indexer.New(cfg,
		indexer.WithRPCPool(pool),
		indexer.WithCheckpointStore(store),
		indexer.WithPostgresSink(pgConfig),
	)
	if err != nil {
		log.Fatalf("Indexer creation failed: %v", err)
	}

	_, err = idx.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start indexer: %v", err)
	}

	fmt.Println("Indexing Uniswap Swaps to PostgreSQL... Waiting for range completion.")

	if err := idx.Wait(); err != nil && err != context.Canceled {
		log.Fatalf("Indexer error: %v", err)
	}

	_ = idx.Stop()
	fmt.Println("Completed indexing range to PostgreSQL!")
}
