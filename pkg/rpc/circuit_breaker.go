package rpc

import (
	"strings"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateHalfOpen:
		return "HALF_OPEN"
	case StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreaker protects RPC endpoints from cascading failures and rate limits.
type CircuitBreaker struct {
	mu sync.Mutex

	state            State
	failureCount     int
	failureThreshold int
	cooldownDuration time.Duration
	lastFailureTime  time.Time
	halfOpenSuccess  int
	successThreshold int
}

// NewCircuitBreaker creates a new CircuitBreaker.
func NewCircuitBreaker(failureThreshold int, cooldownDuration time.Duration) *CircuitBreaker {
	if failureThreshold <= 0 {
		failureThreshold = 5
	}
	if cooldownDuration <= 0 {
		cooldownDuration = 10 * time.Second
	}
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		cooldownDuration: cooldownDuration,
		successThreshold: 2,
		state:            StateClosed,
	}
}

// Allow returns true if a request is permitted to proceed.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	if cb.state == StateOpen {
		if now.Sub(cb.lastFailureTime) > cb.cooldownDuration {
			cb.state = StateHalfOpen
			cb.halfOpenSuccess = 0
			return true
		}
		return false
	}

	return true
}

// RecordResult records the outcome of an RPC call.
func (cb *CircuitBreaker) RecordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err == nil {
		if cb.state == StateHalfOpen {
			cb.halfOpenSuccess++
			if cb.halfOpenSuccess >= cb.successThreshold {
				cb.state = StateClosed
				cb.failureCount = 0
			}
		} else if cb.state == StateClosed {
			cb.failureCount = 0
		}
		return
	}

	// Determine if error is a rate limit (429) or transient server failure
	isRateLimit := IsRateLimitError(err)
	isServerErr := IsServerError(err)

	if isRateLimit || isServerErr {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		if isRateLimit || cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
		}
	}
}

// State returns current state.
func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// IsRateLimitError checks whether an error is due to HTTP 429 / rate limits.
func IsRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "429") ||
		strings.Contains(s, "too many requests") ||
		strings.Contains(s, "rate limit") ||
		strings.Contains(s, "request limit exceeded") ||
		strings.Contains(s, "throughput limit")
}

// IsServerError checks whether an error is a 5xx or connection error.
func IsServerError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "500") ||
		strings.Contains(s, "502") ||
		strings.Contains(s, "503") ||
		strings.Contains(s, "504") ||
		strings.Contains(s, "bad gateway") ||
		strings.Contains(s, "service unavailable") ||
		strings.Contains(s, "connection refused") ||
		strings.Contains(s, "reset by peer")
}
