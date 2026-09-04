package indexer

import (
	"fmt"

	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/metrics"
	"github.com/gNimit/geth-indexer/pkg/rpc"
	"github.com/gNimit/geth-indexer/pkg/sink"
)

// Option represents a functional configuration option for Indexer.
type Option func(*Indexer) error

// WithRPCPool configures a custom RPC pool.
func WithRPCPool(pool *rpc.Pool) Option {
	return func(idx *Indexer) error {
		if pool == nil {
			return fmt.Errorf("indexer: rpc pool cannot be nil")
		}
		idx.pool = pool
		return nil
	}
}

// WithCheckpointStore configures a custom state store.
func WithCheckpointStore(store checkpoint.Store) Option {
	return func(idx *Indexer) error {
		if store == nil {
			return fmt.Errorf("indexer: checkpoint store cannot be nil")
		}
		idx.store = store
		return nil
	}
}

// WithMetrics configures metrics collection.
func WithMetrics(collector metrics.Collector) Option {
	return func(idx *Indexer) error {
		if collector == nil {
			collector = metrics.NewNopCollector()
		}
		idx.metrics = collector
		return nil
	}
}

// WithSink registers an event consumer sink.
func WithSink(s sink.Sink) Option {
	return func(idx *Indexer) error {
		if s == nil {
			return fmt.Errorf("indexer: sink cannot be nil")
		}
		idx.sinks = append(idx.sinks, s)
		return nil
	}
}

// WithPostgresSink configures a PostgreSQL batch sink.
func WithPostgresSink(cfg sink.PostgresConfig) Option {
	return func(idx *Indexer) error {
		pgSink, err := sink.NewPostgresSink(cfg)
		if err != nil {
			return fmt.Errorf("indexer: failed to initialize postgres sink: %w", err)
		}
		idx.sinks = append(idx.sinks, pgSink)
		return nil
	}
}

// WithChannelSink configures the internal channel sink.
func WithChannelSink(opts sink.ChannelSinkOptions) Option {
	return func(idx *Indexer) error {
		idx.channelSink = sink.NewChannelSink(opts)
		idx.sinks = append(idx.sinks, idx.channelSink)
		return nil
	}
}
