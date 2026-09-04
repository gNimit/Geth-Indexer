package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusCollector implements Collector using Prometheus.
type PrometheusCollector struct {
	reg *prometheus.Registry

	eventsProcessed *prometheus.CounterVec
	currentBlock    *prometheus.GaugeVec
	rpcLatency      *prometheus.HistogramVec
	rpcRequests     *prometheus.CounterVec
	batchProcessing *prometheus.HistogramVec
	sinkEvents      *prometheus.CounterVec
	chunkSize       prometheus.Gauge

	mu sync.Mutex
}

// NewPrometheusCollector creates a new PrometheusCollector registered with custom metrics.
func NewPrometheusCollector(namespace string) *PrometheusCollector {
	if namespace == "" {
		namespace = "geth_indexer"
	}

	reg := prometheus.NewRegistry()

	c := &PrometheusCollector{
		reg: reg,
		eventsProcessed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "events_processed_total",
				Help:      "Total number of Ethereum contract events successfully indexed.",
			},
			[]string{"status"},
		),
		currentBlock: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "current_indexed_block",
				Help:      "The highest block number indexed so far.",
			},
			[]string{"chain_id"},
		),
		rpcLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "rpc_latency_seconds",
				Help:      "Histogram of RPC call durations.",
				Buckets:   []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"provider", "method"},
		),
		rpcRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rpc_requests_total",
				Help:      "Total number of RPC requests dispatched to providers.",
			},
			[]string{"provider", "method", "status"},
		),
		batchProcessing: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "sink_batch_processing_seconds",
				Help:      "Duration of sinking batches to destination systems.",
				Buckets:   []float64{0.005, 0.01, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5},
			},
			[]string{"sink"},
		),
		sinkEvents: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "sink_events_total",
				Help:      "Total number of events sent to sinks.",
			},
			[]string{"sink", "status"},
		),
		chunkSize: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "adaptive_chunk_size_blocks",
				Help:      "Current block chunk size dynamically adapted by the fetcher.",
			},
		),
	}

	reg.MustRegister(
		c.eventsProcessed,
		c.currentBlock,
		c.rpcLatency,
		c.rpcRequests,
		c.batchProcessing,
		c.sinkEvents,
		c.chunkSize,
	)

	return c
}

func (p *PrometheusCollector) IncEventsProcessed(count int) {
	if p == nil {
		return
	}
	p.eventsProcessed.WithLabelValues("success").Add(float64(count))
}

func (p *PrometheusCollector) SetCurrentBlock(chainID uint64, block uint64) {
	if p == nil {
		return
	}
	p.currentBlock.WithLabelValues(strconv.FormatUint(chainID, 10)).Set(float64(block))
}

func (p *PrometheusCollector) ObserveRPCLatency(provider string, method string, duration time.Duration) {
	if p == nil {
		return
	}
	p.rpcLatency.WithLabelValues(provider, method).Observe(duration.Seconds())
}

func (p *PrometheusCollector) IncRPCRequests(provider string, method string, status string) {
	if p == nil {
		return
	}
	p.rpcRequests.WithLabelValues(provider, method, status).Inc()
}

func (p *PrometheusCollector) ObserveBatchProcessing(sinkName string, duration time.Duration) {
	if p == nil {
		return
	}
	p.batchProcessing.WithLabelValues(sinkName).Observe(duration.Seconds())
}

func (p *PrometheusCollector) IncSinkEvents(sinkName string, count int, status string) {
	if p == nil {
		return
	}
	p.sinkEvents.WithLabelValues(sinkName, status).Add(float64(count))
}

func (p *PrometheusCollector) SetChunkSize(size uint64) {
	if p == nil {
		return
	}
	p.chunkSize.Set(float64(size))
}

func (p *PrometheusCollector) Handler() http.Handler {
	if p == nil {
		return http.NotFoundHandler()
	}
	return promhttp.HandlerFor(p.reg, promhttp.HandlerOpts{})
}
