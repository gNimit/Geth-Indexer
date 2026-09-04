package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gNimit/geth-indexer/pkg/metrics"
)

// Pool manages a collection of RPC endpoints for fault-tolerant, load-balanced RPC operations.
type Pool struct {
	chainID   uint64
	endpoints []*Endpoint
	balancer  Balancer
	metrics   metrics.Collector

	healthStop   chan struct{}
	healthTicker *time.Ticker

	closeOnce sync.Once
	mu        sync.RWMutex
}

// PoolConfig specifies parameters for initializing an RPCPool.
type PoolConfig struct {
	ChainID        uint64
	Endpoints      []EndpointConfig
	Balancer       Balancer
	Metrics        metrics.Collector
	HealthInterval time.Duration
}

// NewPool creates an initialized Pool.
func NewPool(cfg PoolConfig, customClients []EthClient) (*Pool, error) {
	if len(cfg.Endpoints) == 0 && len(customClients) == 0 {
		return nil, errors.New("rpc: at least one endpoint configuration is required")
	}

	balancer := cfg.Balancer
	if balancer == nil {
		balancer = NewLeastInFlightBalancer()
	}

	metricCollector := cfg.Metrics
	if metricCollector == nil {
		metricCollector = metrics.NewNopCollector()
	}

	var endpoints []*Endpoint
	for i, epCfg := range cfg.Endpoints {
		var custom EthClient
		if i < len(customClients) {
			custom = customClients[i]
		}
		ep, err := NewEndpoint(epCfg, custom)
		if err != nil {
			return nil, fmt.Errorf("rpc: failed to initialize endpoint %s: %w", epCfg.URL, err)
		}
		endpoints = append(endpoints, ep)
	}

	// Handle case where custom clients were passed without explicit EndpointConfigs
	if len(cfg.Endpoints) == 0 && len(customClients) > 0 {
		for i, custom := range customClients {
			epCfg := EndpointConfig{
				URL:      fmt.Sprintf("mock-client-%d", i),
				Provider: "mock",
				Weight:   1,
			}
			ep, err := NewEndpoint(epCfg, custom)
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, ep)
		}
	}

	p := &Pool{
		chainID:    cfg.ChainID,
		endpoints:  endpoints,
		balancer:   balancer,
		metrics:    metricCollector,
		healthStop: make(chan struct{}),
	}

	healthInterval := cfg.HealthInterval
	if healthInterval <= 0 {
		healthInterval = 30 * time.Second
	}
	p.healthTicker = time.NewTicker(healthInterval)
	go p.healthCheckLoop()

	return p, nil
}

// Endpoints returns a slice copy of current endpoints.
func (p *Pool) Endpoints() []*Endpoint {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]*Endpoint, len(p.endpoints))
	copy(out, p.endpoints)
	return out
}

// Execute executes an RPC operation with automatic load balancing, rate limiting, and failover retries.
func (p *Pool) execute(ctx context.Context, method string, fn func(client EthClient) error) error {
	p.mu.RLock()
	endpoints := p.endpoints
	p.mu.RUnlock()

	maxAttempts := len(endpoints)
	if maxAttempts > 3 {
		maxAttempts = 3
	}

	var lastErr error
	tried := make(map[*Endpoint]bool, maxAttempts)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Filter out endpoints already tried in this call
		candidates := make([]*Endpoint, 0, len(endpoints))
		for _, ep := range endpoints {
			if !tried[ep] && ep.IsAvailable() {
				candidates = append(candidates, ep)
			}
		}

		if len(candidates) == 0 {
			if lastErr != nil {
				return fmt.Errorf("rpc pool exhausted retries, last error: %w", lastErr)
			}
			return ErrNoAvailableEndpoints
		}

		ep, err := p.balancer.Select(candidates)
		if err != nil {
			return err
		}
		tried[ep] = true

		if err := ep.Acquire(ctx); err != nil {
			return err
		}

		start := time.Now()
		callErr := fn(ep.Client())
		duration := time.Since(start)

		ep.Release(duration, callErr)
		p.metrics.ObserveRPCLatency(ep.Config().Provider, method, duration)

		if callErr == nil {
			p.metrics.IncRPCRequests(ep.Config().Provider, method, "200")
			return nil
		}

		p.metrics.IncRPCRequests(ep.Config().Provider, method, "error")
		lastErr = callErr

		// If context was canceled, do not retry
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	return fmt.Errorf("rpc call failed after %d attempts: %w", maxAttempts, lastErr)
}

// FilterLogs queries contract logs with load balancing and failover across endpoints.
func (p *Pool) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	var logs []types.Log
	err := p.execute(ctx, "eth_getLogs", func(client EthClient) error {
		res, err := client.FilterLogs(ctx, q)
		if err != nil {
			return err
		}
		logs = res
		return nil
	})
	return logs, err
}

// BlockNumber returns the latest block height across healthy endpoints.
func (p *Pool) BlockNumber(ctx context.Context) (uint64, error) {
	var block uint64
	err := p.execute(ctx, "eth_blockNumber", func(client EthClient) error {
		num, err := client.BlockNumber(ctx)
		if err != nil {
			return err
		}
		block = num
		return nil
	})
	return block, err
}

// HeaderByNumber returns the block header for a given number.
func (p *Pool) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var header *types.Header
	err := p.execute(ctx, "eth_getHeaderByNumber", func(client EthClient) error {
		res, err := client.HeaderByNumber(ctx, number)
		if err != nil {
			return err
		}
		header = res
		return nil
	})
	return header, err
}

// CallContract executes a contract view call using load balancing.
func (p *Pool) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var data []byte
	err := p.execute(ctx, "eth_call", func(client EthClient) error {
		res, err := client.CallContract(ctx, msg, blockNumber)
		if err != nil {
			return err
		}
		data = res
		return nil
	})
	return data, err
}

// healthCheckLoop periodically probes all endpoints for connectivity and block progression.
func (p *Pool) healthCheckLoop() {
	for {
		select {
		case <-p.healthStop:
			return
		case <-p.healthTicker.C:
			p.checkEndpoints()
		}
	}
}

func (p *Pool) checkEndpoints() {
	p.mu.RLock()
	endpoints := p.endpoints
	p.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, ep := range endpoints {
		start := time.Now()
		block, err := ep.Client().BlockNumber(ctx)
		dur := time.Since(start)

		if err != nil {
			ep.SetHealthy(false)
		} else {
			ep.SetHealthy(true)
			ep.SetLastBlock(block)
			ep.Release(dur, nil)
		}
	}
}

// Close stops background health routines and closes all underlying RPC connections.
func (p *Pool) Close() {
	p.closeOnce.Do(func() {
		close(p.healthStop)
		if p.healthTicker != nil {
			p.healthTicker.Stop()
		}

		p.mu.Lock()
		defer p.mu.Unlock()
		for _, ep := range p.endpoints {
			ep.Close()
		}
	})
}
