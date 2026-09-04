package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/rpc"
)

func TestMultiChainIndexer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parsedABI, err := abi.JSON(strings.NewReader(testERC20ABI))
	if err != nil {
		t.Fatal(err)
	}
	transferTopic := parsedABI.Events["Transfer"].ID

	contractEth := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")  // USDC on Ethereum
	contractPoly := common.HexToAddress("0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174") // USDC on Polygon

	mockEthClient := &mockRPCClient{
		latestBlock: 50,
		logsFn: func(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
			return []ethTypes.Log{
				{
					Address:     contractEth,
					Topics:      []common.Hash{transferTopic, common.Hash{}, common.Hash{}},
					Data:        common.LeftPadBytes(big.NewInt(100).Bytes(), 32),
					BlockNumber: q.FromBlock.Uint64(),
					TxHash:      common.HexToHash("0xeth1"),
					Index:       0,
				},
			}, nil
		},
	}

	mockPolyClient := &mockRPCClient{
		latestBlock: 50,
		logsFn: func(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
			return []ethTypes.Log{
				{
					Address:     contractPoly,
					Topics:      []common.Hash{transferTopic, common.Hash{}, common.Hash{}},
					Data:        common.LeftPadBytes(big.NewInt(200).Bytes(), 32),
					BlockNumber: q.FromBlock.Uint64(),
					TxHash:      common.HexToHash("0xpoly1"),
					Index:       0,
				},
			}, nil
		},
	}

	ethPool, _ := rpc.NewPool(rpc.PoolConfig{ChainID: 1, HealthInterval: time.Hour}, []rpc.EthClient{mockEthClient})
	defer ethPool.Close()

	polyPool, _ := rpc.NewPool(rpc.PoolConfig{ChainID: 137, HealthInterval: time.Hour}, []rpc.EthClient{mockPolyClient})
	defer polyPool.Close()

	sharedStore := checkpoint.NewMemoryStore()

	multiCfg := MultiChainConfig{
		Chains: []ChainConfig{
			{
				Name: "ethereum",
				Config: Config{
					ChainID:        1,
					StartBlock:     1,
					EndBlock:       5,
					BlockChunkSize: 5,
					Contracts: []ContractConfig{
						{Address: contractEth, ABI: &parsedABI, Events: []string{"Transfer"}},
					},
				},
				RPCPool: WithRPCPool(ethPool),
			},
			{
				Name: "polygon",
				Config: Config{
					ChainID:        137,
					StartBlock:     1,
					EndBlock:       5,
					BlockChunkSize: 5,
					Contracts: []ContractConfig{
						{Address: contractPoly, ABI: &parsedABI, Events: []string{"Transfer"}},
					},
				},
				RPCPool: WithRPCPool(polyPool),
			},
		},
	}

	multiIndexer, err := NewMultiChain(multiCfg, WithMultiChainStore(sharedStore))
	if err != nil {
		t.Fatalf("failed to initialize MultiChainIndexer: %v", err)
	}

	unifiedCh, err := multiIndexer.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start MultiChainIndexer: %v", err)
	}

	receivedByChain := make(map[uint64]int)
	for ev := range unifiedCh {
		receivedByChain[ev.ChainID]++
		if receivedByChain[1] >= 1 && receivedByChain[137] >= 1 {
			_ = multiIndexer.Stop()
			break
		}
	}

	if receivedByChain[1] == 0 || receivedByChain[137] == 0 {
		t.Fatalf("expected events from both chain 1 and chain 137, got: %+v", receivedByChain)
	}
	fmt.Printf("Received events across multiple chains: %+v\n", receivedByChain)
}
