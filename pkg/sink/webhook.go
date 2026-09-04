package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gNimit/geth-indexer/pkg/types"
)

// WebhookConfig configures HTTP webhook streaming.
type WebhookConfig struct {
	URL     string            `json:"url" yaml:"url"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Timeout time.Duration     `json:"timeout" yaml:"timeout"`
	Retries int               `json:"retries" yaml:"retries"`
}

// WebhookSink forwards event batches to an external HTTP webhook endpoint.
type WebhookSink struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewWebhookSink creates a new WebhookSink.
func NewWebhookSink(cfg WebhookConfig) *WebhookSink {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	retries := cfg.Retries
	if retries <= 0 {
		retries = 3
	}
	cfg.Retries = retries

	return &WebhookSink{
		cfg: cfg,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (w *WebhookSink) Name() string {
	return "webhook"
}

func (w *WebhookSink) Send(ctx context.Context, batch []*types.Event) error {
	if len(batch) == 0 {
		return nil
	}

	payload, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("webhook: failed to marshal batch: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < w.cfg.Retries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.cfg.URL, bytes.NewReader(payload))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		for k, v := range w.cfg.Headers {
			req.Header.Set(k, v)
		}

		resp, err := w.client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			lastErr = fmt.Errorf("webhook: non-2xx status received: %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(100*(1<<attempt)) * time.Millisecond):
		}
	}

	return fmt.Errorf("webhook: failed after %d retries: %w", w.cfg.Retries, lastErr)
}

func (w *WebhookSink) Close() error {
	w.client.CloseIdleConnections()
	return nil
}

// StdoutSink logs events directly to standard output or logger.
type StdoutSink struct {
	JSONFormat bool
}

// NewStdoutSink creates a StdoutSink.
func NewStdoutSink(jsonFormat bool) *StdoutSink {
	return &StdoutSink{JSONFormat: jsonFormat}
}

func (s *StdoutSink) Name() string {
	return "stdout"
}

func (s *StdoutSink) Send(ctx context.Context, batch []*types.Event) error {
	for _, ev := range batch {
		if s.JSONFormat {
			b, _ := json.Marshal(ev)
			fmt.Println(string(b))
		} else {
			fmt.Println(ev.String())
		}
	}
	return nil
}

func (s *StdoutSink) Close() error {
	return nil
}
