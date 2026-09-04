package metrics

import (
	"net/http"
	"time"
)

// Collector defines the interface for monitoring geth-indexer metrics.
type Collector interface {
	// IncEventsProcessed increments the total number of indexed events.
	IncEventsProcessed(count int)

	// SetCurrentBlock updates the current indexed block height gauge.
	SetCurrentBlock(chainID uint64, block uint64)

	// ObserveRPCLatency records the duration of an RPC call.
	ObserveRPCLatency(provider string, method string, duration time.Duration)

	// IncRPCRequests increments the count of RPC requests made.
	IncRPCRequests(provider string, method string, status string)

	// ObserveBatchProcessing records the time taken to process and sink a batch.
	ObserveBatchProcessing(sinkName string, duration time.Duration)

	// IncSinkEvents increments events sent to a sink.
	IncSinkEvents(sinkName string, count int, status string)

	// SetChunkSize records the current adaptive block chunk size.
	SetChunkSize(size uint64)

	// Handler returns an http.Handler for exposing metrics (e.g. /metrics).
	Handler() http.Handler
}
