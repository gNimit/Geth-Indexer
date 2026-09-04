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

const erc20ABI = `[
	{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},
	{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}
]`

const uniswapSwapABI = `[
	{"anonymous":false,"inputs":[{"indexed":true,"name":"sender","type":"address"},{"indexed":true,"name":"recipient","type":"address"},{"indexed":false,"name":"amount0","type":"int256"},{"indexed":false,"name":"amount1","type":"int256"}],"name":"Swap","type":"event"}
]`

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 1. Ethereum Mainnet RPC Pool (ChainID 1)
	ethPool, err := rpc.NewPool(rpc.PoolConfig{
		ChainID: 1,
		Endpoints: []rpc.EndpointConfig{
			{URL: "https://eth.llamarpc.com", Provider: "llama-eth", RPS: 20},
		},
	}, nil)
	if err != nil {
		log.Fatalf("Failed to create Ethereum RPC pool: %v", err)
	}

	// 2. Arbitrum One RPC Pool (ChainID 42161)
	arbPool, err := rpc.NewPool(rpc.PoolConfig{
		ChainID: 42161,
		Endpoints: []rpc.EndpointConfig{
			{URL: "https://arb1.arbitrum.io/rpc", Provider: "arb-public", RPS: 20},
		},
	}, nil)
	if err != nil {
		log.Fatalf("Failed to create Arbitrum RPC pool: %v", err)
	}

	// 3. Checkpoint Store
	store, _ := checkpoint.NewFileStore("multichain_checkpoint.json")

	// 4. Configure Multiple Chains with Multiple Contracts and Multiple Events
	multiConfig := indexer.MultiChainConfig{
		Chains: []indexer.ChainConfig{
			{
				Name: "ethereum",
				Config: indexer.Config{
					ChainID:        1,
					StartBlock:     19500000,
					EndBlock:       0, // live tailing
					BlockChunkSize: 500,
					Concurrency:    4,
					Confirmations:  12,
					PollInterval:   2 * time.Second,
					Contracts: []indexer.ContractConfig{
						// Contract 1: USDC on Ethereum (Transfer + Approval)
						{
							Address: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
							ABIJSON: erc20ABI,
							Events:  []string{"Transfer", "Approval"},
						},
						// Contract 2: Uniswap V3 Pool on Ethereum (Swap)
						{
							Address: common.HexToAddress("0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"),
							ABIJSON: uniswapSwapABI,
							Events:  []string{"Swap"},
						},
					},
				},
				RPCPool: indexer.WithRPCPool(ethPool),
			},
			{
				Name: "arbitrum",
				Config: indexer.Config{
					ChainID:        42161,
					StartBlock:     190000000,
					EndBlock:       0, // live tailing
					BlockChunkSize: 1000,
					Concurrency:    4,
					Confirmations:  20,
					PollInterval:   1 * time.Second,
					Contracts: []indexer.ContractConfig{
						// Contract 3: ARB Token on Arbitrum (Transfer)
						{
							Address: common.HexToAddress("0x912CE59144191C1204E64559FE8253a0e49E6548"),
							ABIJSON: erc20ABI,
							Events:  []string{"Transfer"},
						},
					},
				},
				RPCPool: indexer.WithRPCPool(arbPool),
			},
		},
	}

	multiIndexer, err := indexer.NewMultiChain(multiConfig,
		indexer.WithMultiChainStore(store),
	)
	if err != nil {
		log.Fatalf("Failed to initialize multi-chain indexer: %v", err)
	}

	// 5. Start multi-chain indexing and receive unified channel
	unifiedEventsCh, err := multiIndexer.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start multi-chain indexer: %v", err)
	}

	fmt.Println("Listening for events across Ethereum and Arbitrum... Press Ctrl+C to stop.")

	for event := range unifiedEventsCh {
		fmt.Printf("[Chain %d] %s at %s | Block: %d | Tx: %s | Data: %v\n",
			event.ChainID,
			event.EventName,
			event.ContractAddress.Hex(),
			event.BlockNumber,
			event.TxHash.Hex(),
			event.Data,
		)
	}

	_ = multiIndexer.Stop()
	fmt.Println("Multi-chain indexer stopped.")
}
