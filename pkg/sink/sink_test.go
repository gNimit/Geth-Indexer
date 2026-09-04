package sink

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gNimit/geth-indexer/pkg/types"
)

func TestChannelSink(t *testing.T) {
	ctx := context.Background()
	cs := NewChannelSink(ChannelSinkOptions{BufferSize: 10, Blocking: true})
	defer cs.Close()

	if cs.Name() != "channel" {
		t.Fatalf("expected name 'channel', got '%s'", cs.Name())
	}

	events := []*types.Event{
		{
			ChainID:         1,
			ContractAddress: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			EventName:       "Transfer",
			BlockNumber:     100,
			TxHash:          common.HexToHash("0x123"),
			LogIndex:        0,
		},
		{
			ChainID:         1,
			ContractAddress: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			EventName:       "Transfer",
			BlockNumber:     100,
			TxHash:          common.HexToHash("0x123"),
			LogIndex:        1,
		},
	}

	if err := cs.Send(ctx, events); err != nil {
		t.Fatalf("send failed: %v", err)
	}

	ch := cs.Events()
	ev1 := <-ch
	ev2 := <-ch

	if ev1.LogIndex != 0 || ev2.LogIndex != 1 {
		t.Fatalf("received unexpected events: %v, %v", ev1, ev2)
	}
}

func TestWebhookSink(t *testing.T) {
	var receivedCount int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&receivedCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	sink := NewWebhookSink(WebhookConfig{
		URL:     server.URL,
		Timeout: 2 * time.Second,
		Retries: 2,
	})
	defer sink.Close()

	events := []*types.Event{
		{
			ChainID:     1,
			EventName:   "Transfer",
			BlockNumber: 100,
		},
	}

	if err := sink.Send(context.Background(), events); err != nil {
		t.Fatalf("webhook send failed: %v", err)
	}

	if atomic.LoadInt64(&receivedCount) != 1 {
		t.Fatalf("expected 1 webhook request, got %d", receivedCount)
	}
}

func TestStdoutSink(t *testing.T) {
	s := NewStdoutSink(false)
	events := []*types.Event{
		{
			ChainID:         1,
			ContractAddress: common.HexToAddress("0x1234567890123456789012345678901234567890"),
			EventName:       "Transfer",
			BlockNumber:     100,
			TxHash:          common.HexToHash("0xabc"),
			LogIndex:        0,
		},
	}
	if err := s.Send(context.Background(), events); err != nil {
		t.Fatalf("stdout send failed: %v", err)
	}
}
