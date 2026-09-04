package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPrometheusCollector(t *testing.T) {
	c := NewPrometheusCollector("test_indexer")
	c.IncEventsProcessed(10)
	c.SetCurrentBlock(1, 19000000)
	c.ObserveRPCLatency("infura", "eth_getLogs", 50*time.Millisecond)
	c.IncRPCRequests("infura", "eth_getLogs", "200")
	c.ObserveBatchProcessing("postgres", 20*time.Millisecond)
	c.IncSinkEvents("postgres", 10, "success")
	c.SetChunkSize(1000)

	handler := c.Handler()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if len(body) == 0 {
		t.Fatal("expected non-empty metrics output")
	}
}

func TestNopCollector(t *testing.T) {
	c := NewNopCollector()
	c.IncEventsProcessed(10)
	c.SetCurrentBlock(1, 100)
	c.ObserveRPCLatency("alchemy", "eth_getLogs", time.Second)
	c.IncRPCRequests("alchemy", "eth_getLogs", "200")
	c.ObserveBatchProcessing("channel", time.Millisecond)
	c.IncSinkEvents("channel", 5, "success")
	c.SetChunkSize(500)

	handler := c.Handler()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}
