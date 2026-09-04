package checkpoint

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// MemoryStore is an in-memory thread-safe implementation of Store useful for tests or ephemeral tasks.
type MemoryStore struct {
	checkpoints     map[string]uint64
	unindexedRanges map[string][]BlockRange
	mu              sync.RWMutex
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		checkpoints:     make(map[string]uint64),
		unindexedRanges: make(map[string][]BlockRange),
	}
}

func (m *MemoryStore) GetLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address) (uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.checkpoints[stateKey(chainID, contract)], nil
}

func (m *MemoryStore) SaveLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address, block uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := stateKey(chainID, contract)
	if block > m.checkpoints[key] {
		m.checkpoints[key] = block
	}
	return nil
}

func (m *MemoryStore) GetUnindexedRanges(ctx context.Context, chainID uint64, contract common.Address) ([]BlockRange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := stateKey(chainID, contract)
	var pending []BlockRange
	for _, r := range m.unindexedRanges[key] {
		if !r.Completed {
			pending = append(pending, r)
		}
	}
	return pending, nil
}

func (m *MemoryStore) AddUnindexedRange(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := stateKey(chainID, contract)
	r.UpdatedAt = time.Now()
	m.unindexedRanges[key] = append(m.unindexedRanges[key], r)
	return nil
}

func (m *MemoryStore) MarkRangeCompleted(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := stateKey(chainID, contract)
	for i := range m.unindexedRanges[key] {
		if m.unindexedRanges[key][i].FromBlock == r.FromBlock && m.unindexedRanges[key][i].ToBlock == r.ToBlock {
			m.unindexedRanges[key][i].Completed = true
			m.unindexedRanges[key][i].UpdatedAt = time.Now()
			break
		}
	}
	return nil
}

func (m *MemoryStore) Close() error {
	return nil
}
