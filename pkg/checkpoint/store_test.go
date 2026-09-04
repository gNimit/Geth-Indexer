package checkpoint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func testStoreContract(t *testing.T, store Store) {
	ctx := context.Background()
	chainID := uint64(1)
	contract := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

	// Initially zero
	val, err := store.GetLastIndexedBlock(ctx, chainID, contract)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 0 {
		t.Fatalf("expected 0, got %d", val)
	}

	// Save progress
	if err := store.SaveLastIndexedBlock(ctx, chainID, contract, 100); err != nil {
		t.Fatalf("failed to save checkpoint: %v", err)
	}

	val, err = store.GetLastIndexedBlock(ctx, chainID, contract)
	if err != nil || val != 100 {
		t.Fatalf("expected 100, got %d (err: %v)", val, err)
	}

	// Save lower block shouldn't regress
	if err := store.SaveLastIndexedBlock(ctx, chainID, contract, 90); err != nil {
		t.Fatalf("failed to save: %v", err)
	}
	val, _ = store.GetLastIndexedBlock(ctx, chainID, contract)
	if val != 100 {
		t.Fatalf("checkpoint regressed to %d, expected 100", val)
	}

	// Unindexed ranges tracking
	r1 := BlockRange{FromBlock: 50, ToBlock: 60}
	r2 := BlockRange{FromBlock: 70, ToBlock: 80}
	if err := store.AddUnindexedRange(ctx, chainID, contract, r1); err != nil {
		t.Fatalf("failed to add range: %v", err)
	}
	if err := store.AddUnindexedRange(ctx, chainID, contract, r2); err != nil {
		t.Fatalf("failed to add range: %v", err)
	}

	ranges, err := store.GetUnindexedRanges(ctx, chainID, contract)
	if err != nil || len(ranges) != 2 {
		t.Fatalf("expected 2 unindexed ranges, got %d (err: %v)", len(ranges), err)
	}

	// Mark r1 completed
	if err := store.MarkRangeCompleted(ctx, chainID, contract, r1); err != nil {
		t.Fatalf("failed to mark range completed: %v", err)
	}

	ranges, err = store.GetUnindexedRanges(ctx, chainID, contract)
	if err != nil || len(ranges) != 1 || ranges[0].FromBlock != 70 {
		t.Fatalf("expected 1 pending range [70, 80], got %+v", ranges)
	}
}

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	testStoreContract(t, store)
}

func TestFileStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "checkpoint_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "checkpoint.json")
	store, err := NewFileStore(filePath)
	if err != nil {
		t.Fatalf("failed to create file store: %v", err)
	}

	testStoreContract(t, store)

	// Re-open from disk to verify persistence across restarts
	reopened, err := NewFileStore(filePath)
	if err != nil {
		t.Fatalf("failed to reload file store: %v", err)
	}

	ctx := context.Background()
	contract := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	val, err := reopened.GetLastIndexedBlock(ctx, 1, contract)
	if err != nil || val != 100 {
		t.Fatalf("expected 100 after reload, got %d (err: %v)", val, err)
	}

	ranges, err := reopened.GetUnindexedRanges(ctx, 1, contract)
	if err != nil || len(ranges) != 1 || ranges[0].FromBlock != 70 {
		t.Fatalf("expected 1 pending range after reload, got %+v", ranges)
	}
}
