package rpc

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
)

// Balancer defines the strategy for selecting an available RPC endpoint.
type Balancer interface {
	Select(endpoints []*Endpoint) (*Endpoint, error)
	Name() string
}

// RoundRobinBalancer selects endpoints sequentially.
type RoundRobinBalancer struct {
	counter uint64
}

// NewRoundRobinBalancer creates a round-robin load balancer.
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{}
}

func (b *RoundRobinBalancer) Name() string {
	return "round_robin"
}

func (b *RoundRobinBalancer) Select(endpoints []*Endpoint) (*Endpoint, error) {
	n := len(endpoints)
	if n == 0 {
		return nil, ErrNoAvailableEndpoints
	}

	startIdx := atomic.AddUint64(&b.counter, 1) - 1
	for i := 0; i < n; i++ {
		idx := (startIdx + uint64(i)) % uint64(n)
		ep := endpoints[idx]
		if ep.IsAvailable() {
			return ep, nil
		}
	}
	return nil, ErrNoAvailableEndpoints
}

// LeastInFlightBalancer routes requests to the endpoint with the fewest active in-flight requests.
type LeastInFlightBalancer struct{}

// NewLeastInFlightBalancer creates a least-in-flight load balancer.
func NewLeastInFlightBalancer() *LeastInFlightBalancer {
	return &LeastInFlightBalancer{}
}

func (b *LeastInFlightBalancer) Name() string {
	return "least_in_flight"
}

func (b *LeastInFlightBalancer) Select(endpoints []*Endpoint) (*Endpoint, error) {
	if len(endpoints) == 0 {
		return nil, ErrNoAvailableEndpoints
	}

	var best *Endpoint
	minInFlight := int64(math.MaxInt64)

	for _, ep := range endpoints {
		if !ep.IsAvailable() {
			continue
		}
		inFlight := ep.InFlight()
		if inFlight < minInFlight {
			minInFlight = inFlight
			best = ep
		}
	}

	if best == nil {
		return nil, ErrNoAvailableEndpoints
	}
	return best, nil
}

// LatencyWeightedBalancer selects the fastest responsive endpoint with weighted random exploration.
type LatencyWeightedBalancer struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// NewLatencyWeightedBalancer creates a latency-weighted balancer.
func NewLatencyWeightedBalancer() *LatencyWeightedBalancer {
	return &LatencyWeightedBalancer{
		rng: rand.New(rand.NewSource(rand.Int63())),
	}
}

func (b *LatencyWeightedBalancer) Name() string {
	return "latency_weighted"
}

func (b *LatencyWeightedBalancer) Select(endpoints []*Endpoint) (*Endpoint, error) {
	if len(endpoints) == 0 {
		return nil, ErrNoAvailableEndpoints
	}

	var available []*Endpoint
	for _, ep := range endpoints {
		if ep.IsAvailable() {
			available = append(available, ep)
		}
	}

	if len(available) == 0 {
		return nil, ErrNoAvailableEndpoints
	}
	if len(available) == 1 {
		return available[0], nil
	}

	// Calculate inverse latency weights (lower latency = higher weight)
	weights := make([]float64, len(available))
	totalWeight := 0.0

	for i, ep := range available {
		lat := ep.LatencyEMA()
		if lat <= 0 {
			lat = 10.0 // default optimistic estimate
		}
		// Weight is inversely proportional to latency, scaled by endpoint weight config
		w := (1000.0 / lat) * float64(ep.Config().Weight)
		weights[i] = w
		totalWeight += w
	}

	b.mu.Lock()
	r := b.rng.Float64() * totalWeight
	b.mu.Unlock()

	cur := 0.0
	for i, w := range weights {
		cur += w
		if r <= cur {
			return available[i], nil
		}
	}

	return available[0], nil
}
