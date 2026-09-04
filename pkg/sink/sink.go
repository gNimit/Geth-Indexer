package sink

import (
	"context"

	"github.com/gNimit/geth-indexer/pkg/types"
)

// Sink defines the pluggable interface for consuming indexed event batches.
// Consumers like Kafka, gRPC, PostgreSQL, WebSockets, or Go channels implement this.
type Sink interface {
	// Name returns a human-readable identifier for logging and metrics.
	Name() string

	// Send dispatches a batch of events to the destination system.
	Send(ctx context.Context, batch []*types.Event) error

	// Close flushes buffers and cleans up connections.
	Close() error
}
