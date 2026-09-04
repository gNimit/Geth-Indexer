package sink

import (
	"context"
	"errors"
	"sync"

	"github.com/gNimit/geth-indexer/pkg/types"
)

var ErrChannelClosed = errors.New("sink: channel is closed")

// ChannelSink streams events into a native Go channel with deadlock prevention.
type ChannelSink struct {
	ch       chan *types.Event
	stopCh   chan struct{}
	blocking bool
	mu       sync.RWMutex
	closed   bool
}

// ChannelSinkOptions configures ChannelSink behavior.
type ChannelSinkOptions struct {
	BufferSize int  // Buffer size of the internal Go channel
	Blocking   bool // If true, Send will wait for buffer space; if false, non-blocking
}

// NewChannelSink creates an initialized ChannelSink.
func NewChannelSink(opts ChannelSinkOptions) *ChannelSink {
	bufSize := opts.BufferSize
	if bufSize <= 0 {
		bufSize = 20000
	}
	return &ChannelSink{
		ch:       make(chan *types.Event, bufSize),
		stopCh:   make(chan struct{}),
		blocking: opts.Blocking,
	}
}

func (s *ChannelSink) Name() string {
	return "channel"
}

// Events returns the read-only Go channel of events.
func (s *ChannelSink) Events() <-chan *types.Event {
	return s.ch
}

// Send emits the batch of events to the channel, respecting context and stop signals.
func (s *ChannelSink) Send(ctx context.Context, batch []*types.Event) error {
	s.mu.RLock()
	closed := s.closed
	s.mu.RUnlock()

	if closed {
		return ErrChannelClosed
	}

	for _, event := range batch {
		if s.blocking {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-s.stopCh:
				return ErrChannelClosed
			case s.ch <- event:
			}
		} else {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-s.stopCh:
				return ErrChannelClosed
			case s.ch <- event:
			default:
				// Non-blocking drop or wait with context
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-s.stopCh:
					return ErrChannelClosed
				case s.ch <- event:
				}
			}
		}
	}
	return nil
}

// Close safely terminates the channel sink and unblocks any waiting senders.
func (s *ChannelSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.closed {
		s.closed = true
		close(s.stopCh)
		close(s.ch)
	}
	return nil
}
