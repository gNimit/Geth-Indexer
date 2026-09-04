package indexer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/metrics"
	"github.com/gNimit/geth-indexer/pkg/sink"
	"github.com/gNimit/geth-indexer/pkg/types"
)

// ChainConfig defines the indexing configuration for a single blockchain network.
type ChainConfig struct {
	Name    string // e.g. "ethereum", "arbitrum", "polygon", "base"
	Config  Config // Core indexer configuration (Contracts, StartBlock, EndBlock, etc.)
	RPCPool Option // withRPCPool option specifically for this chain
	Store   checkpoint.Store
}

// MultiChainConfig defines configuration across multiple EVM chains.
type MultiChainConfig struct {
	Chains            []ChainConfig
	ChannelBufferSize int
}

// MultiChainIndexer coordinates independent Indexer instances across multiple chains
// and merges their output into a single unified event channel and pluggable sinks.
type MultiChainIndexer struct {
	indexers    map[uint64]*Indexer
	channelSink *sink.ChannelSink
	sinks       []sink.Sink
	metrics     metrics.Collector
	store       checkpoint.Store

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	errCh  chan error

	mu      sync.Mutex
	started bool
	stopped bool
}

// MultiChainOption provides configuration options for MultiChainIndexer.
type MultiChainOption func(*MultiChainIndexer)

// WithMultiChainMetrics sets metrics collector across all chains.
func WithMultiChainMetrics(collector metrics.Collector) MultiChainOption {
	return func(m *MultiChainIndexer) {
		m.metrics = collector
	}
}

// WithMultiChainSink registers a sink that receives events from all chains.
func WithMultiChainSink(s sink.Sink) MultiChainOption {
	return func(m *MultiChainIndexer) {
		m.sinks = append(m.sinks, s)
	}
}

// WithMultiChainStore sets a shared checkpoint store across all chains.
func WithMultiChainStore(s checkpoint.Store) MultiChainOption {
	return func(m *MultiChainIndexer) {
		m.store = s
	}
}

// NewMultiChain creates a new MultiChainIndexer managing multiple chains.
func NewMultiChain(cfg MultiChainConfig, opts ...MultiChainOption) (*MultiChainIndexer, error) {
	if len(cfg.Chains) == 0 {
		return nil, errors.New("multichain: at least one chain configuration is required")
	}

	bufSize := cfg.ChannelBufferSize
	if bufSize <= 0 {
		bufSize = 50000
	}

	chSink := sink.NewChannelSink(sink.ChannelSinkOptions{
		BufferSize: bufSize,
		Blocking:   false,
	})

	m := &MultiChainIndexer{
		indexers:    make(map[uint64]*Indexer),
		channelSink: chSink,
		errCh:       make(chan error, len(cfg.Chains)),
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.metrics == nil {
		m.metrics = metrics.NewNopCollector()
	}
	if m.store == nil {
		m.store = checkpoint.NewMemoryStore()
	}

	// Initialize individual indexers for each chain
	for _, cc := range cfg.Chains {
		chainID := cc.Config.ChainID
		if chainID == 0 {
			return nil, fmt.Errorf("multichain: chain %s must specify a non-zero ChainID", cc.Name)
		}
		if _, exists := m.indexers[chainID]; exists {
			return nil, fmt.Errorf("multichain: duplicate ChainID %d", chainID)
		}

		chainStore := cc.Store
		if chainStore == nil {
			chainStore = m.store
		}

		chainOpts := []Option{
			WithCheckpointStore(chainStore),
			WithMetrics(m.metrics),
			WithSink(chSink),
		}
		if cc.RPCPool != nil {
			chainOpts = append(chainOpts, cc.RPCPool)
		}
		for _, s := range m.sinks {
			chainOpts = append(chainOpts, WithSink(s))
		}

		idx, err := New(cc.Config, chainOpts...)
		if err != nil {
			return nil, fmt.Errorf("multichain: failed to create indexer for chain %s (%d): %w", cc.Name, chainID, err)
		}
		m.indexers[chainID] = idx
	}

	return m, nil
}

// Start launches all chain indexers concurrently and returns the unified event channel.
func (m *MultiChainIndexer) Start(ctx context.Context) (<-chan *types.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return nil, errors.New("multichain: already started")
	}
	m.started = true
	m.ctx, m.cancel = context.WithCancel(ctx)

	for chainID, idx := range m.indexers {
		cID := chainID
		cIdx := idx
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			_, err := cIdx.Start(m.ctx)
			if err != nil {
				m.errCh <- fmt.Errorf("chain %d error: %w", cID, err)
				return
			}
			if err := cIdx.Wait(); err != nil && !errors.Is(err, context.Canceled) {
				m.errCh <- fmt.Errorf("chain %d terminated with error: %w", cID, err)
			}
		}()
	}

	return m.channelSink.Events(), nil
}

// Stop terminates all chain indexers gracefully.
func (m *MultiChainIndexer) Stop() error {
	m.mu.Lock()
	if !m.started || m.stopped {
		m.mu.Unlock()
		return nil
	}
	m.stopped = true
	m.cancel()
	m.mu.Unlock()

	var errs []error
	for chainID, idx := range m.indexers {
		if err := idx.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("chain %d stop error: %w", chainID, err))
		}
	}

	m.wg.Wait()
	_ = m.channelSink.Close()

	if len(errs) > 0 {
		return fmt.Errorf("multichain stop errors: %v", errs)
	}
	return nil
}

// Wait blocks until all chain indexers finish or any encounters an unrecoverable error.
func (m *MultiChainIndexer) Wait() error {
	select {
	case <-m.ctx.Done():
		return m.ctx.Err()
	case err := <-m.errCh:
		return err
	}
}

// ChainIndexer returns the specific Indexer instance for a chain.
func (m *MultiChainIndexer) ChainIndexer(chainID uint64) (*Indexer, bool) {
	idx, ok := m.indexers[chainID]
	return idx, ok
}
