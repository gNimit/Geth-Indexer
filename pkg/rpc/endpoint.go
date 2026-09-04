package rpc

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/time/rate"
)

// EthClient defines the interface required for Ethereum node interaction.
// *ethclient.Client naturally satisfies this interface.
type EthClient interface {
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	Close()
}

// EndpointConfig defines configuration for an individual RPC endpoint.
type EndpointConfig struct {
	URL              string        `json:"url" yaml:"url"`
	Provider         string        `json:"provider" yaml:"provider"` // e.g. "alchemy", "infura", "local"
	Weight           int           `json:"weight" yaml:"weight"`     // relative weight for balancing (default 1)
	RPS              float64       `json:"rps" yaml:"rps"`           // max requests per second (0 = unlimited)
	Burst            int           `json:"burst" yaml:"burst"`       // max burst requests
	Timeout          time.Duration `json:"timeout" yaml:"timeout"`   // per-request timeout
	MaxRetries       int           `json:"maxRetries" yaml:"maxRetries"`
	FailureThreshold int           `json:"failureThreshold" yaml:"failureThreshold"`
	CooldownDuration time.Duration `json:"cooldownDuration" yaml:"cooldownDuration"`
}

// Endpoint wraps a connection to a specific RPC provider with telemetry, rate limiting, and circuit breaking.
type Endpoint struct {
	config EndpointConfig
	client EthClient

	limiter        *rate.Limiter
	circuitBreaker *CircuitBreaker

	inFlight    int64
	latencyEMA  float64 // Exponential moving average of latency in milliseconds
	lastBlock   uint64  // Latest observed block from this node
	lastChecked time.Time
	healthy     bool

	mu sync.RWMutex
}

// NewEndpoint creates a managed Endpoint. If customClient is provided, it is used directly.
func NewEndpoint(cfg EndpointConfig, customClient EthClient) (*Endpoint, error) {
	if cfg.Weight <= 0 {
		cfg.Weight = 1
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.CooldownDuration <= 0 {
		cfg.CooldownDuration = 15 * time.Second
	}

	var limiter *rate.Limiter
	if cfg.RPS > 0 {
		burst := cfg.Burst
		if burst <= 0 {
			burst = int(cfg.RPS)
			if burst < 1 {
				burst = 1
			}
		}
		limiter = rate.NewLimiter(rate.Limit(cfg.RPS), burst)
	}

	client := customClient
	if client == nil {
		dialed, err := ethclient.Dial(cfg.URL)
		if err != nil {
			return nil, err
		}
		client = dialed
	}

	ep := &Endpoint{
		config:         cfg,
		client:         client,
		limiter:        limiter,
		circuitBreaker: NewCircuitBreaker(cfg.FailureThreshold, cfg.CooldownDuration),
		healthy:        true,
	}
	return ep, nil
}

// Config returns endpoint configuration.
func (e *Endpoint) Config() EndpointConfig {
	return e.config
}

// IsAvailable checks if the endpoint can accept a request (circuit breaker and health).
func (e *Endpoint) IsAvailable() bool {
	if !e.circuitBreaker.Allow() {
		return false
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.healthy
}

// InFlight returns the number of active requests currently executing on this endpoint.
func (e *Endpoint) InFlight() int64 {
	return atomic.LoadInt64(&e.inFlight)
}

// LatencyEMA returns the exponential moving average latency in milliseconds.
func (e *Endpoint) LatencyEMA() float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.latencyEMA
}

// LastBlock returns the latest block height observed on this endpoint.
func (e *Endpoint) LastBlock() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.lastBlock
}

// SetLastBlock sets the observed block number.
func (e *Endpoint) SetLastBlock(block uint64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lastBlock = block
}

// SetHealthy sets the health status.
func (e *Endpoint) SetHealthy(healthy bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.healthy = healthy
	e.lastChecked = time.Now()
}

// Acquire waits for rate limiter if configured.
func (e *Endpoint) Acquire(ctx context.Context) error {
	if e.limiter != nil {
		if err := e.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	atomic.AddInt64(&e.inFlight, 1)
	return nil
}

// Release records request duration and releases in-flight counter.
func (e *Endpoint) Release(duration time.Duration, err error) {
	atomic.AddInt64(&e.inFlight, -1)
	e.circuitBreaker.RecordResult(err)

	e.mu.Lock()
	defer e.mu.Unlock()
	durMs := float64(duration.Milliseconds())
	if e.latencyEMA == 0 {
		e.latencyEMA = durMs
	} else {
		// Smoothing factor alpha = 0.2
		e.latencyEMA = 0.2*durMs + 0.8*e.latencyEMA
	}
}

// Close closes the underlying client connection.
func (e *Endpoint) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

// Client returns the underlying EthClient.
func (e *Endpoint) Client() EthClient {
	return e.client
}

var ErrNoAvailableEndpoints = errors.New("rpc: no available healthy endpoints in pool")
