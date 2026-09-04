package checkpoint

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// SQLStore implements Store backed by any database/sql compatible driver (e.g., PostgreSQL).
type SQLStore struct {
	db *sql.DB
}

// NewSQLStore creates a SQLStore and ensures checkpoint tables are initialized.
func NewSQLStore(db *sql.DB) (*SQLStore, error) {
	if db == nil {
		return nil, fmt.Errorf("checkpoint: db connection cannot be nil")
	}

	s := &SQLStore{db: db}
	if err := s.initSchema(); err != nil {
		return nil, fmt.Errorf("checkpoint: failed to initialize schema: %w", err)
	}

	return s, nil
}

func (s *SQLStore) initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS geth_indexer_checkpoints (
			chain_id BIGINT NOT NULL,
			contract_address VARCHAR(42) NOT NULL,
			last_indexed_block BIGINT NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
			PRIMARY KEY (chain_id, contract_address)
		);`,
		`CREATE TABLE IF NOT EXISTS geth_indexer_unindexed_ranges (
			chain_id BIGINT NOT NULL,
			contract_address VARCHAR(42) NOT NULL,
			from_block BIGINT NOT NULL,
			to_block BIGINT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
			PRIMARY KEY (chain_id, contract_address, from_block, to_block)
		);`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (s *SQLStore) GetLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address) (uint64, error) {
	query := `SELECT last_indexed_block FROM geth_indexer_checkpoints WHERE chain_id = $1 AND contract_address = $2`
	var lastBlock uint64
	err := s.db.QueryRowContext(ctx, query, chainID, contract.Hex()).Scan(&lastBlock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return lastBlock, nil
}

func (s *SQLStore) SaveLastIndexedBlock(ctx context.Context, chainID uint64, contract common.Address, block uint64) error {
	query := `
		INSERT INTO geth_indexer_checkpoints (chain_id, contract_address, last_indexed_block, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (chain_id, contract_address) DO UPDATE
		SET last_indexed_block = GREATEST(geth_indexer_checkpoints.last_indexed_block, EXCLUDED.last_indexed_block),
		    updated_at = EXCLUDED.updated_at
	`
	_, err := s.db.ExecContext(ctx, query, chainID, contract.Hex(), block, time.Now())
	return err
}

func (s *SQLStore) GetUnindexedRanges(ctx context.Context, chainID uint64, contract common.Address) ([]BlockRange, error) {
	query := `
		SELECT from_block, to_block, completed, updated_at
		FROM geth_indexer_unindexed_ranges
		WHERE chain_id = $1 AND contract_address = $2 AND completed = FALSE
		ORDER BY from_block ASC
	`
	rows, err := s.db.QueryContext(ctx, query, chainID, contract.Hex())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranges []BlockRange
	for rows.Next() {
		var r BlockRange
		if err := rows.Scan(&r.FromBlock, &r.ToBlock, &r.Completed, &r.UpdatedAt); err != nil {
			return nil, err
		}
		ranges = append(ranges, r)
	}
	return ranges, rows.Err()
}

func (s *SQLStore) AddUnindexedRange(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error {
	query := `
		INSERT INTO geth_indexer_unindexed_ranges (chain_id, contract_address, from_block, to_block, completed, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (chain_id, contract_address, from_block, to_block) DO NOTHING
	`
	_, err := s.db.ExecContext(ctx, query, chainID, contract.Hex(), r.FromBlock, r.ToBlock, r.Completed, time.Now())
	return err
}

func (s *SQLStore) MarkRangeCompleted(ctx context.Context, chainID uint64, contract common.Address, r BlockRange) error {
	query := `
		UPDATE geth_indexer_unindexed_ranges
		SET completed = TRUE, updated_at = $1
		WHERE chain_id = $2 AND contract_address = $3 AND from_block = $4 AND to_block = $5
	`
	_, err := s.db.ExecContext(ctx, query, time.Now(), chainID, contract.Hex(), r.FromBlock, r.ToBlock)
	return err
}

func (s *SQLStore) Close() error {
	return nil // connection lifecycle managed externally
}
