package checkpoint

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// BlockRange represents an interval of blocks to index or track.
type BlockRange struct {
	FromBlock uint64    `json:"fromBlock"`
	ToBlock   uint64    `json:"toBlock"`
	Completed bool      `json:"completed"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Store defines the state tracking interface for indexer checkpoints and gap recovery.
type Store interface {
	// GetLastIndexedBlock returns the highest successfully indexed block for a contract.
	GetLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address) (uint64, error)

	// SaveLastIndexedBlock records progress up to a given block number.
	SaveLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address, block uint64) error

	// GetUnindexedRanges returns any historical gaps or ranges pending processing.
	GetUnindexedRanges(ctx context.Context, chainID uint64, contract common.Address) ([]BlockRange, error)

	// AddUnindexedRange records a block interval that needs indexing.
	AddUnindexedRange(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error

	// MarkRangeCompleted marks an unindexed range as resolved.
	MarkRangeCompleted(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error

	// Close safely flushes and closes underlying resources.
	Close() error
}
