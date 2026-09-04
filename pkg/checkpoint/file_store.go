package checkpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type contractState struct {
	LastIndexedBlock uint64       `json:"lastIndexedBlock"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	UnindexedRanges  []BlockRange `json:"unindexedRanges"`
}

type fileState struct {
	Contracts map[string]*contractState `json:"contracts"` // key: "chainID:address"
}

// FileStore implements Store with an atomically-written local JSON file.
type FileStore struct {
	filePath string
	state    fileState
	mu       sync.RWMutex
}

// NewFileStore initializes or loads state from a JSON file.
func NewFileStore(filePath string) (*FileStore, error) {
	if filePath == "" {
		filePath = "indexer_checkpoint.json"
	}

	fs := &FileStore{
		filePath: filePath,
		state: fileState{
			Contracts: make(map[string]*contractState),
		},
	}

	if err := fs.load(); err != nil {
		return nil, err
	}

	return fs, nil
}

func stateKey(chainID uint64, contract common.Address) string {
	return fmt.Sprintf("%d:%s", chainID, contract.Hex())
}

func (fs *FileStore) load() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("checkpoint: failed to read file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &fs.state); err != nil {
		return fmt.Errorf("checkpoint: failed to parse JSON state: %w", err)
	}
	if fs.state.Contracts == nil {
		fs.state.Contracts = make(map[string]*contractState)
	}

	return nil
}

// persist writes state atomically using a temporary file and atomic rename.
func (fs *FileStore) persist() error {
	dir := filepath.Dir(fs.filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("checkpoint: failed to create directory: %w", err)
		}
	}

	data, err := json.MarshalIndent(fs.state, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: failed to serialize state: %w", err)
	}

	tmpPath := fmt.Sprintf("%s.tmp.%d", fs.filePath, time.Now().UnixNano())
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("checkpoint: failed to write temp file: %w", err)
	}

	if err := os.Rename(tmpPath, fs.filePath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("checkpoint: failed to commit state atomically: %w", err)
	}

	return nil
}

func (fs *FileStore) GetLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address) (uint64, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	key := stateKey(chainID, contract)
	if cs, ok := fs.state.Contracts[key]; ok {
		return cs.LastIndexedBlock, nil
	}
	return 0, nil
}

func (fs *FileStore) SaveLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address, block uint64) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	key := stateKey(chainID, contract)
	cs, ok := fs.state.Contracts[key]
	if !ok {
		cs = &contractState{}
		fs.state.Contracts[key] = cs
	}

	if block > cs.LastIndexedBlock {
		cs.LastIndexedBlock = block
	}
	cs.UpdatedAt = time.Now()

	return fs.persist()
}

func (fs *FileStore) GetUnindexedRanges(ctx context.Context, chainID uint64, contract common.Address) ([]BlockRange, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	key := stateKey(chainID, contract)
	cs, ok := fs.state.Contracts[key]
	if !ok {
		return nil, nil
	}

	var pending []BlockRange
	for _, r := range cs.UnindexedRanges {
		if !r.Completed {
			pending = append(pending, r)
		}
	}
	return pending, nil
}

func (fs *FileStore) AddUnindexedRange(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	key := stateKey(chainID, contract)
	cs, ok := fs.state.Contracts[key]
	if !ok {
		cs = &contractState{}
		fs.state.Contracts[key] = cs
	}

	r.UpdatedAt = time.Now()
	cs.UnindexedRanges = append(cs.UnindexedRanges, r)

	return fs.persist()
}

func (fs *FileStore) MarkRangeCompleted(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	key := stateKey(chainID, contract)
	cs, ok := fs.state.Contracts[key]
	if !ok {
		return nil
	}

	for i := range cs.UnindexedRanges {
		if cs.UnindexedRanges[i].FromBlock == r.FromBlock && cs.UnindexedRanges[i].ToBlock == r.ToBlock {
			cs.UnindexedRanges[i].Completed = true
			cs.UnindexedRanges[i].UpdatedAt = time.Now()
			break
		}
	}

	return fs.persist()
}

func (fs *FileStore) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.persist()
}
