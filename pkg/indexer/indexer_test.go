package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/metrics"
	"github.com/gNimit/geth-indexer/pkg/rpc"
)

const testERC20ABI = `[
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

type mockRPCClient struct {
	latestBlock uint64
	logsFn      func(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error)
}

func (m *mockRPCClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
	if m.logsFn != nil {
		return m.logsFn(ctx, q)
	}
	return []ethTypes.Log{}, nil
}

func (m *mockRPCClient) BlockNumber(ctx context.Context) (uint64, error) {
	return m.latestBlock, nil
}

func (m *mockRPCClient) HeaderByNumber(ctx context.Context, number *big.Int) (*ethTypes.Header, error) {
	return &ethTypes.Header{Number: number}, nil
}

func (m *mockRPCClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}

func (m *mockRPCClient) Close() {}

func TestIndexerEndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	contractAddr := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	parsedABI, err := abi.JSON(strings.NewReader(testERC20ABI))
	if err != nil {
		t.Fatal(err)
	}
	transferTopic := parsedABI.Events["Transfer"].ID

	fromAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	toAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")

	mockClient := &mockRPCClient{
		latestBlock: 200,
		logsFn: func(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
			from := q.FromBlock.Uint64()
			to := q.ToBlock.Uint64()

			var logs []ethTypes.Log
			for b := from; b <= to; b++ {
				logs = append(logs,
					ethTypes.Log{
						Address:     contractAddr,
						Topics:      []common.Hash{transferTopic, common.BytesToHash(fromAddr.Bytes()), common.BytesToHash(toAddr.Bytes())},
						Data:        common.LeftPadBytes(big.NewInt(int64(b*10)).Bytes(), 32),
						BlockNumber: b,
						TxHash:      common.HexToHash(fmt.Sprintf("0x%x", b)),
						Index:       0,
					},
					ethTypes.Log{
						Address:     contractAddr,
						Topics:      []common.Hash{transferTopic, common.BytesToHash(fromAddr.Bytes()), common.BytesToHash(toAddr.Bytes())},
						Data:        common.LeftPadBytes(big.NewInt(int64(b*10+1)).Bytes(), 32),
						BlockNumber: b,
						TxHash:      common.HexToHash(fmt.Sprintf("0x%x", b)),
						Index:       1,
					},
				)
			}
			return logs, nil
		},
	}

	pool, err := rpc.NewPool(rpc.PoolConfig{
		ChainID:        1,
		HealthInterval: time.Hour,
	}, []rpc.EthClient{mockClient})
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	store := checkpoint.NewMemoryStore()
	collector := metrics.NewPrometheusCollector("test_e2e")

	cfg := Config{
		ChainID:    1,
		StartBlock: 1,
		EndBlock:   50, // 50 blocks * 2 events = 100 events
		Contracts: []ContractConfig{
			{
				Address: contractAddr,
				ABI:     &parsedABI,
				Events:  []string{"Transfer"},
			},
		},
		BlockChunkSize:    10,
		Concurrency:       4,
		Confirmations:     0,
		SinkBatchSize:     20,
		SinkFlushInterval: 10 * time.Millisecond,
	}

	idx, err := New(cfg,
		WithRPCPool(pool),
		WithCheckpointStore(store),
		WithMetrics(collector),
	)
	if err != nil {
		t.Fatalf("failed to create indexer: %v", err)
	}

	eventsCh, err := idx.Start(ctx)
	if err != nil {
		t.Fatalf("failed to start indexer: %v", err)
	}

	receivedCount := 0
	var lastBlock uint64 = 0
	var lastTxIndex uint = 0
	var lastLogIndex uint = 0

	for ev := range eventsCh {
		receivedCount++

		if ev.BlockNumber < lastBlock {
			t.Fatalf("ordering regression: got block %d after %d", ev.BlockNumber, lastBlock)
		}
		if ev.BlockNumber == lastBlock && ev.TxIndex == lastTxIndex && ev.LogIndex < lastLogIndex {
			t.Fatalf("log index regression within block %d: %d < %d", ev.BlockNumber, ev.LogIndex, lastLogIndex)
		}

		lastBlock = ev.BlockNumber
		lastTxIndex = ev.TxIndex
		lastLogIndex = ev.LogIndex

		if receivedCount == 100 {
			_ = idx.Stop()
			break
		}
	}

	if receivedCount != 100 {
		t.Fatalf("expected 100 events, got %d", receivedCount)
	}

	lastSaved, err := store.GetLastIndexedBlock(ctx, 1, contractAddr)
	if err != nil || lastSaved < 50 {
		t.Fatalf("expected checkpoint >= 50, got %d (err: %v)", lastSaved, err)
	}
}

func TestAdaptiveChunkBisection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	contractAddr := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	parsedABI, _ := abi.JSON(strings.NewReader(testERC20ABI))
	transferTopic := parsedABI.Events["Transfer"].ID

	var calls int64
	mockClient := &mockRPCClient{
		latestBlock: 100,
		logsFn: func(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
			atomic.AddInt64(&calls, 1)
			from := q.FromBlock.Uint64()
			to := q.ToBlock.Uint64()

			if to-from > 5 {
				return nil, fmt.Errorf("query returned more than 10000 results")
			}

			return []ethTypes.Log{
				{
					Address:     contractAddr,
					Topics:      []common.Hash{transferTopic, common.Hash{}, common.Hash{}},
					Data:        common.LeftPadBytes(big.NewInt(1).Bytes(), 32),
					BlockNumber: from,
					TxHash:      common.HexToHash("0x1"),
					Index:       0,
				},
			}, nil
		},
	}

	pool, _ := rpc.NewPool(rpc.PoolConfig{ChainID: 1, HealthInterval: time.Hour}, []rpc.EthClient{mockClient})
	defer pool.Close()
	store := checkpoint.NewMemoryStore()

	cfg := Config{
		ChainID:           1,
		StartBlock:        1,
		EndBlock:          20,
		BlockChunkSize:    20, // starts at 20, which exceeds 5, forcing bisection
		Concurrency:       1,
		Confirmations:     0,
		SinkBatchSize:     10,
		SinkFlushInterval: 10 * time.Millisecond,
		Contracts: []ContractConfig{
			{Address: contractAddr, ABI: &parsedABI},
		},
	}

	idx, err := New(cfg, WithRPCPool(pool), WithCheckpointStore(store))
	if err != nil {
		t.Fatal(err)
	}

	eventsCh, err := idx.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}

	eventsReceived := 0
	for range eventsCh {
		eventsReceived++
		if eventsReceived >= 4 {
			_ = idx.Stop()
			break
		}
	}

	if eventsReceived == 0 {
		t.Fatal("expected events after bisection, got none")
	}

	if atomic.LoadInt64(&calls) <= 1 {
		t.Fatalf("expected multiple calls due to range bisection, got %d", calls)
	}
}

func BenchmarkIndexerThroughput(b *testing.B) {
	contractAddr := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	parsedABI, _ := abi.JSON(strings.NewReader(testERC20ABI))
	transferTopic := parsedABI.Events["Transfer"].ID

	logsPerBlock := 100
	mockClient := &mockRPCClient{
		latestBlock: 500000,
		logsFn: func(ctx context.Context, q ethereum.FilterQuery) ([]ethTypes.Log, error) {
			from := q.FromBlock.Uint64()
			to := q.ToBlock.Uint64()
			numBlocks := int(to - from + 1)
			totalLogs := numBlocks * logsPerBlock
			logs := make([]ethTypes.Log, totalLogs)

			idx := 0
			for blk := from; blk <= to; blk++ {
				for i := 0; i < logsPerBlock; i++ {
					logs[idx] = ethTypes.Log{
						Address:     contractAddr,
						Topics:      []common.Hash{transferTopic, common.Hash{}, common.Hash{}},
						Data:        common.LeftPadBytes(big.NewInt(int64(i)).Bytes(), 32),
						BlockNumber: blk,
						TxHash:      common.HexToHash("0xabc"),
						Index:       uint(i),
					}
					idx++
				}
			}
			return logs, nil
		},
	}

	// Benchmark processing 2,000 events per iteration
	const eventsPerIter = 2000
	blocksNeeded := uint64(eventsPerIter / logsPerBlock)

	cfg := Config{
		ChainID:           1,
		StartBlock:        1,
		EndBlock:          blocksNeeded,
		BlockChunkSize:    10,
		Concurrency:       4,
		Confirmations:     0,
		SinkBatchSize:     500,
		SinkFlushInterval: 5 * time.Millisecond,
		ChannelBufferSize: 5000,
		Contracts: []ContractConfig{
			{Address: contractAddr, ABI: &parsedABI},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pool, _ := rpc.NewPool(rpc.PoolConfig{ChainID: 1, HealthInterval: time.Hour}, []rpc.EthClient{mockClient})
		store := checkpoint.NewMemoryStore()

		idx, err := New(cfg, WithRPCPool(pool), WithCheckpointStore(store))
		if err != nil {
			b.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		eventsCh, err := idx.Start(ctx)
		if err != nil {
			cancel()
			b.Fatal(err)
		}

		received := 0
		for range eventsCh {
			received++
			if received >= eventsPerIter {
				cancel()
				_ = idx.Stop()
				break
			}
		}
	}
}
