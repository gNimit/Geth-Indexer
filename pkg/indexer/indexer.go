package indexer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gNimit/geth-indexer/internal/abi"
	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/metrics"
	"github.com/gNimit/geth-indexer/pkg/rpc"
	"github.com/gNimit/geth-indexer/pkg/sink"
	"github.com/gNimit/geth-indexer/pkg/types"
)

// Indexer is the core engine coordinating RPC fetching, ABI decoding, checkpointing, and sink dispatching.
type Indexer struct {
	config      Config
	pool        *rpc.Pool
	registry    *abi.Registry
	store       checkpoint.Store
	metrics     metrics.Collector
	sinks       []sink.Sink
	channelSink *sink.ChannelSink

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	errCh  chan error

	mu      sync.Mutex
	started bool
	stopped bool
}

// New creates and configures an Indexer instance with functional options.
func New(cfg Config, opts ...Option) (*Indexer, error) {
	defaults := DefaultConfig()
	if cfg.BlockChunkSize == 0 {
		cfg.BlockChunkSize = defaults.BlockChunkSize
	}
	if cfg.MinBlockChunkSize == 0 {
		cfg.MinBlockChunkSize = defaults.MinBlockChunkSize
	}
	if cfg.MaxBlockChunkSize == 0 {
		cfg.MaxBlockChunkSize = defaults.MaxBlockChunkSize
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = defaults.Concurrency
	}
	if cfg.Confirmations == 0 {
		cfg.Confirmations = defaults.Confirmations
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = defaults.PollInterval
	}
	if cfg.ChannelBufferSize <= 0 {
		cfg.ChannelBufferSize = defaults.ChannelBufferSize
	}
	if cfg.SinkBatchSize <= 0 {
		cfg.SinkBatchSize = defaults.SinkBatchSize
	}
	if cfg.SinkFlushInterval <= 0 {
		cfg.SinkFlushInterval = defaults.SinkFlushInterval
	}

	registry := abi.NewRegistry()
	for _, c := range cfg.Contracts {
		if c.ABI != nil {
			if err := registry.RegisterParsedContract(c.Address, c.ABI, c.Events); err != nil {
				return nil, fmt.Errorf("indexer: failed to register parsed contract %s: %w", c.Address.Hex(), err)
			}
		} else if c.ABIJSON != "" {
			if err := registry.RegisterContract(c.Address, c.ABIJSON, c.Events); err != nil {
				return nil, fmt.Errorf("indexer: failed to parse ABI for %s: %w", c.Address.Hex(), err)
			}
		}
	}

	idx := &Indexer{
		config:   cfg,
		registry: registry,
		errCh:    make(chan error, 1),
	}

	for _, opt := range opts {
		if err := opt(idx); err != nil {
			return nil, err
		}
	}

	// Ensure default ChannelSink exists so callers can stream directly from Start()
	if idx.channelSink == nil {
		idx.channelSink = sink.NewChannelSink(sink.ChannelSinkOptions{
			BufferSize: cfg.ChannelBufferSize,
			Blocking:   false, // Non-blocking with buffer to prevent pipeline stalling
		})
		idx.sinks = append(idx.sinks, idx.channelSink)
	}

	// Default fallback for checkpoint store: MemoryStore
	if idx.store == nil {
		idx.store = checkpoint.NewMemoryStore()
	}

	// Default fallback for metrics: NopCollector
	if idx.metrics == nil {
		idx.metrics = metrics.NewNopCollector()
	}

	return idx, nil
}

// RegisterSink adds a new event consumer sink to the indexer.
func (idx *Indexer) RegisterSink(s sink.Sink) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.started {
		return errors.New("indexer: cannot register sink after indexing has started")
	}
	if s == nil {
		return errors.New("indexer: sink cannot be nil")
	}
	idx.sinks = append(idx.sinks, s)
	return nil
}

// Start launches the indexer pipeline in background goroutines and returns the event channel.
func (idx *Indexer) Start(ctx context.Context) (<-chan *types.Event, error) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idx.started {
		return nil, errors.New("indexer: already started")
	}
	if idx.pool == nil {
		return nil, errors.New("indexer: rpc pool must be configured before starting")
	}

	idx.started = true
	idx.ctx, idx.cancel = context.WithCancel(ctx)

	idx.wg.Add(1)
	go func() {
		defer idx.wg.Done()
		if err := idx.runPipeline(idx.ctx); err != nil && !errors.Is(err, context.Canceled) {
			select {
			case idx.errCh <- err:
			default:
			}
		}
	}()

	return idx.channelSink.Events(), nil
}

// Stop initiates graceful shutdown, flushes all sinks, commits checkpoints, and releases resources.
func (idx *Indexer) Stop() error {
	idx.mu.Lock()
	if !idx.started || idx.stopped {
		idx.mu.Unlock()
		return nil
	}
	idx.stopped = true
	idx.cancel()
	idx.mu.Unlock()

	// Wait for pipeline workers to complete with safety timeout
	done := make(chan struct{})
	go func() {
		idx.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
	}

	var errs []error
	for _, s := range idx.sinks {
		if err := s.Close(); err != nil {
			errs = append(errs, fmt.Errorf("sink %s close error: %w", s.Name(), err))
		}
	}

	if idx.store != nil {
		if err := idx.store.Close(); err != nil {
			errs = append(errs, fmt.Errorf("checkpoint store close error: %w", err))
		}
	}

	if idx.pool != nil {
		idx.pool.Close()
	}

	if len(errs) > 0 {
		return fmt.Errorf("indexer stop encountered errors: %v", errs)
	}
	return nil
}

// Wait blocks until the indexer finishes its configured range or encounters a fatal error.
func (idx *Indexer) Wait() error {
	select {
	case <-idx.ctx.Done():
		return idx.ctx.Err()
	case err := <-idx.errCh:
		return err
	}
}

// ContractAddresses returns the addresses of all monitored contracts.
func (idx *Indexer) ContractAddresses() []common.Address {
	addrs := make([]common.Address, len(idx.config.Contracts))
	for i, c := range idx.config.Contracts {
		addrs[i] = c.Address
	}
	return addrs
}
