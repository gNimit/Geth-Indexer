package metrics

import (
	"net/http"
	"time"
)

// NopCollector is a zero-allocation no-op implementation of Collector.
type NopCollector struct{}

// NewNopCollector creates a new no-op collector.
func NewNopCollector() *NopCollector {
	return &NopCollector{}
}

func (n *NopCollector) IncEventsProcessed(count int)                                            {}
func (n *NopCollector) SetCurrentBlock(chainID uint64, block uint64)                           {}
func (n *NopCollector) ObserveRPCLatency(provider string, method string, duration time.Duration) {}
func (n *NopCollector) IncRPCRequests(provider string, method string, status string)           {}
func (n *NopCollector) ObserveBatchProcessing(sinkName string, duration time.Duration)         {}
func (n *NopCollector) IncSinkEvents(sinkName string, count int, status string)                {}
func (n *NopCollector) SetChunkSize(size uint64)                                               {}
func (n *NopCollector) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("metrics disabled"))
	})
}
