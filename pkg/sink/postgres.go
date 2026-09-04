package sink

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gNimit/geth-indexer/pkg/types"
	_ "github.com/lib/pq"
)

// PostgresConfig configures the PostgreSQL batch sink.
type PostgresConfig struct {
	Host            string        `json:"host" yaml:"host"`
	Port            int           `json:"port" yaml:"port"`
	User            string        `json:"user" yaml:"user"`
	Password        string        `json:"password" yaml:"password"`
	Database        string        `json:"database" yaml:"database"`
	SSLMode         string        `json:"sslMode" yaml:"sslMode"`
	TableName       string        `json:"tableName" yaml:"tableName"`
	MaxOpenConns    int           `json:"maxOpenConns" yaml:"maxOpenConns"`
	MaxIdleConns    int           `json:"maxIdleConns" yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" yaml:"connMaxLifetime"`
}

// PostgresSink batches events and persists them to PostgreSQL idempotently.
type PostgresSink struct {
	db        *sql.DB
	tableName string
}

// NewPostgresSink initializes the PostgreSQL connection, sets connection pools, and ensures the table schema.
func NewPostgresSink(cfg PostgresConfig) (*PostgresSink, error) {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	if cfg.TableName == "" {
		cfg.TableName = "indexed_events"
	}
	if cfg.Port == 0 {
		cfg.Port = 5432
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres sink: connection failed: %w", err)
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 25
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 10
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	sink := &PostgresSink{
		db:        db,
		tableName: cfg.TableName,
	}

	if err := sink.initSchema(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres sink: schema migration failed: %w", err)
	}

	return sink, nil
}

// NewPostgresSinkWithDB creates a PostgresSink using an existing *sql.DB connection.
func NewPostgresSinkWithDB(db *sql.DB, tableName string) (*PostgresSink, error) {
	if db == nil {
		return nil, fmt.Errorf("postgres sink: db cannot be nil")
	}
	if tableName == "" {
		tableName = "indexed_events"
	}
	sink := &PostgresSink{
		db:        db,
		tableName: tableName,
	}
	if err := sink.initSchema(); err != nil {
		return nil, fmt.Errorf("postgres sink: schema migration failed: %w", err)
	}
	return sink, nil
}

func (s *PostgresSink) initSchema() error {
	queries := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			id VARCHAR(100) PRIMARY KEY,
			chain_id BIGINT NOT NULL,
			contract_address VARCHAR(42) NOT NULL,
			event_name VARCHAR(100) NOT NULL,
			event_signature VARCHAR(66) NOT NULL,
			block_number BIGINT NOT NULL,
			block_hash VARCHAR(66) NOT NULL,
			tx_hash VARCHAR(66) NOT NULL,
			tx_index INTEGER NOT NULL,
			log_index INTEGER NOT NULL,
			data JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);`, s.tableName),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_contract_block ON %s (contract_address, block_number);`, s.tableName, s.tableName),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_name_block ON %s (event_name, block_number);`, s.tableName, s.tableName),
		fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_tx_hash ON %s (tx_hash);`, s.tableName, s.tableName),
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostgresSink) Name() string {
	return "postgres"
}

// Send executes a high-speed multi-row parameterized batch INSERT ON CONFLICT DO NOTHING.
func (s *PostgresSink) Send(ctx context.Context, batch []*types.Event) error {
	if len(batch) == 0 {
		return nil
	}

	// PostgreSQL supports up to 65535 query parameters.
	// Each event has 11 parameters: id, chain_id, contract_address, event_name, event_signature,
	// block_number, block_hash, tx_hash, tx_index, log_index, data
	const paramsPerEvent = 11
	const maxBatchSize = 5000 // safe threshold well under 65535 / 11

	for start := 0; start < len(batch); start += maxBatchSize {
		end := start + maxBatchSize
		if end > len(batch) {
			end = len(batch)
		}
		subBatch := batch[start:end]

		if err := s.insertSubBatch(ctx, subBatch, paramsPerEvent); err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresSink) insertSubBatch(ctx context.Context, batch []*types.Event, paramsPerEvent int) error {
	var valPlaceholders []string
	vals := make([]interface{}, 0, len(batch)*paramsPerEvent)

	for i, ev := range batch {
		baseParam := i * paramsPerEvent
		valPlaceholders = append(valPlaceholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			baseParam+1, baseParam+2, baseParam+3, baseParam+4, baseParam+5,
			baseParam+6, baseParam+7, baseParam+8, baseParam+9, baseParam+10, baseParam+11,
		))

		jsonData, err := json.Marshal(ev.Data)
		if err != nil {
			jsonData = []byte("{}")
		}

		vals = append(vals,
			ev.ID(),
			ev.ChainID,
			ev.ContractAddress.Hex(),
			ev.EventName,
			ev.EventSignature.Hex(),
			ev.BlockNumber,
			ev.BlockHash.Hex(),
			ev.TxHash.Hex(),
			ev.TxIndex,
			ev.LogIndex,
			string(jsonData),
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, chain_id, contract_address, event_name, event_signature,
			block_number, block_hash, tx_hash, tx_index, log_index, data
		) VALUES %s
		ON CONFLICT (id) DO NOTHING
	`, s.tableName, strings.Join(valPlaceholders, ","))

	_, err := s.db.ExecContext(ctx, query, vals...)
	if err != nil {
		return fmt.Errorf("postgres sink: batch insert failed: %w", err)
	}
	return nil
}

func (s *PostgresSink) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
