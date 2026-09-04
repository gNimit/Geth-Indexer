package rpc

import (
	"context"
	"errors"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gNimit/geth-indexer/pkg/metrics"
)

type mockEthClient struct {
	filterLogsHook func(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	blockNumHook   func(ctx context.Context) (uint64, error)
	callCount      int64
	closed         bool
}

func (m *mockEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	atomic.AddInt64(&m.callCount, 1)
	if m.filterLogsHook != nil {
		return m.filterLogsHook(ctx, q)
	}
	return []types.Log{}, nil
}

func (m *mockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	atomic.AddInt64(&m.callCount, 1)
	if m.blockNumHook != nil {
		return m.blockNumHook(ctx)
	}
	return 1000, nil
}

func (m *mockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(1000)}, nil
}

func (m *mockEthClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return []byte("ok"), nil
}

func (m *mockEthClient) Close() {
	m.closed = true
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	if !cb.Allow() {
		t.Fatal("circuit breaker should allow on closed state")
	}

	// First normal error (not 429)
	cb.RecordResult(errors.New("random err"))
	if cb.State() != StateClosed {
		t.Fatal("cb should remain closed")
	}

	// 429 rate limit trips circuit breaker immediately
	cb.RecordResult(errors.New("status code: 429 Too Many Requests"))
	if cb.State() != StateOpen {
		t.Fatalf("cb should be OPEN after 429, got %s", cb.State())
	}

	if cb.Allow() {
		t.Fatal("cb should reject while OPEN")
	}

	// Wait for cooldown
	time.Sleep(60 * time.Millisecond)
	if !cb.Allow() {
		t.Fatal("cb should allow trial request after cooldown (HALF_OPEN)")
	}

	// Record successes to close
	cb.RecordResult(nil)
	cb.RecordResult(nil)
	if cb.State() != StateClosed {
		t.Fatalf("cb should be CLOSED after trial successes, got %s", cb.State())
	}
}

func TestPoolFailoverAndLoadBalancing(t *testing.T) {
	ctx := context.Background()

	// Client 1 fails with 429
	client1 := &mockEthClient{
		filterLogsHook: func(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
			return nil, errors.New("429 rate limit exceeded")
		},
	}

	// Client 2 succeeds
	expectedLogs := []types.Log{
		{
			Address:     common.HexToAddress("0x1111111111111111111111111111111111111111"),
			BlockNumber: 500,
		},
	}
	client2 := &mockEthClient{
		filterLogsHook: func(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
			return expectedLogs, nil
		},
	}

	poolCfg := PoolConfig{
		ChainID: 1,
		Endpoints: []EndpointConfig{
			{URL: "http://provider1", Provider: "provider1", Weight: 1},
			{URL: "http://provider2", Provider: "provider2", Weight: 1},
		},
		Balancer:       NewRoundRobinBalancer(),
		Metrics:        metrics.NewNopCollector(),
		HealthInterval: 1 * time.Hour, // don't run health check during test
	}

	pool, err := NewPool(poolCfg, []EthClient{client1, client2})
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	logs, err := pool.FilterLogs(ctx, ethereum.FilterQuery{})
	if err != nil {
		t.Fatalf("expected successful failover, got error: %v", err)
	}

	if len(logs) != 1 || logs[0].BlockNumber != 500 {
		t.Fatalf("unexpected logs returned: %+v", logs)
	}

	// Client1 should have been called and tripped, Client2 should have handled the request
	if atomic.LoadInt64(&client1.callCount) == 0 {
		t.Fatal("client1 was not called")
	}
	if atomic.LoadInt64(&client2.callCount) == 0 {
		t.Fatal("client2 was not called after failover")
	}
}

func TestPoolBalancers(t *testing.T) {
	c1 := &mockEthClient{}
	c2 := &mockEthClient{}

	ep1, _ := NewEndpoint(EndpointConfig{URL: "ep1", Weight: 1}, c1)
	ep2, _ := NewEndpoint(EndpointConfig{URL: "ep2", Weight: 1}, c2)
	endpoints := []*Endpoint{ep1, ep2}

	// Test RoundRobin
	rr := NewRoundRobinBalancer()
	sel1, _ := rr.Select(endpoints)
	sel2, _ := rr.Select(endpoints)
	if sel1 == sel2 {
		t.Fatal("round robin should cycle between endpoints")
	}

	// Test LeastInFlight
	atomic.StoreInt64(&ep1.inFlight, 5)
	atomic.StoreInt64(&ep2.inFlight, 1)
	lif := NewLeastInFlightBalancer()
	selLIF, err := lif.Select(endpoints)
	if err != nil || selLIF != ep2 {
		t.Fatalf("least in flight should have picked ep2, got %v", selLIF)
	}

	// Test LatencyWeighted
	lw := NewLatencyWeightedBalancer()
	selLW, err := lw.Select(endpoints)
	if err != nil || selLW == nil {
		t.Fatalf("latency weighted should pick an endpoint, got err: %v", err)
	}
}
