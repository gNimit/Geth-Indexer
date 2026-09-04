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
)

// Standard ERC-20 ABI snippet
const erc20ABI = `[
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "from", "type": "address"},
			{"indexed": true, "name": "to", "type": "address"},
			{"indexed": false, "name": "value", "type": "uint256"}
		],
		"name": "Transfer",
		"type": "event"
	}
]`

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 1. Initialize RPC Pool with rotation and rate limits
	pool, err := rpc.NewPool(rpc.PoolConfig{
		ChainID: 1,
		Endpoints: []rpc.EndpointConfig{
			{URL: "https://eth.llamarpc.com", Provider: "llama", RPS: 20, Weight: 1},
			{URL: "https://rpc.ankr.com/eth", Provider: "ankr", RPS: 20, Weight: 1},
		},
		Balancer: rpc.NewLeastInFlightBalancer(),
	}, nil)
	if err != nil {
		log.Fatalf("Failed to initialize RPC pool: %v", err)
	}

	// 2. Configure State Checkpointing (atomic local file store)
	store, err := checkpoint.NewFileStore("usdc_checkpoint.json")
	if err != nil {
		log.Fatalf("Failed to initialize checkpoint store: %v", err)
	}

	// 3. Configure Indexer for USDC Contract Transfer events
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	cfg := indexer.Config{
		ChainID:           1,
		StartBlock:        19500000,
		EndBlock:          0, // 0 = continuous live streaming mode
		BlockChunkSize:    500,
		Concurrency:       4,
		Confirmations:     6,
		PollInterval:      2 * time.Second,
		ChannelBufferSize: 10000,
		Contracts: []indexer.ContractConfig{
			{
				Address: usdcAddress,
				ABIJSON: erc20ABI,
				Events:  []string{"Transfer"},
			},
		},
	}

	idx, err := indexer.New(cfg,
		indexer.WithRPCPool(pool),
		indexer.WithCheckpointStore(store),
	)
	if err != nil {
		log.Fatalf("Failed to initialize indexer: %v", err)
	}

	// 4. Start Indexer and receive native Go event channel
	eventsCh, err := idx.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start indexer: %v", err)
	}

	fmt.Println("Listening for USDC Transfer events... Press Ctrl+C to stop.")

	// 5. Consume events directly from the channel (pipe to Kafka, WebSocket, gRPC, etc.)
	for event := range eventsCh {
		fmt.Printf("Event: %s | Block: %d | Tx: %s | From: %v | To: %v | Value: %v\n",
			event.EventName,
			event.BlockNumber,
			event.TxHash.Hex(),
			event.Data["from"],
			event.Data["to"],
			event.Data["value"],
		)
	}

	_ = idx.Stop()
	fmt.Println("Indexer stopped cleanly.")
}
