package indexer

import (
	"time"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// ContractConfig configures a specific smart contract for event listening.
type ContractConfig struct {
	Address    common.Address  `json:"address" yaml:"address"`
	ABI        *ethAbi.ABI     `json:"-" yaml:"-"`
	ABIJSON    string          `json:"abiJson" yaml:"abiJson"` // Raw ABI JSON string if ABI is not pre-parsed
	Events     []string        `json:"events" yaml:"events"`   // Event names to filter (empty = all events)
	StartBlock uint64          `json:"startBlock" yaml:"startBlock"`
}

// Config defines the runtime configuration for the Indexer engine.
type Config struct {
	ChainID            uint64           `json:"chainId" yaml:"chainId"`
	Contracts          []ContractConfig `json:"contracts" yaml:"contracts"`
	StartBlock         uint64           `json:"startBlock" yaml:"startBlock"`
	EndBlock           uint64           `json:"endBlock" yaml:"endBlock"` // 0 = continuous live tailing
	BlockChunkSize     uint64           `json:"blockChunkSize" yaml:"blockChunkSize"`
	MinBlockChunkSize  uint64           `json:"minBlockChunkSize" yaml:"minBlockChunkSize"`
	MaxBlockChunkSize  uint64           `json:"maxBlockChunkSize" yaml:"maxBlockChunkSize"`
	Concurrency        int              `json:"concurrency" yaml:"concurrency"`
	Confirmations      uint64           `json:"confirmations" yaml:"confirmations"` // Reorg safety buffer
	PollInterval       time.Duration    `json:"pollInterval" yaml:"pollInterval"`
	ChannelBufferSize  int              `json:"channelBufferSize" yaml:"channelBufferSize"`
	SinkBatchSize      int              `json:"sinkBatchSize" yaml:"sinkBatchSize"`
	SinkFlushInterval  time.Duration    `json:"sinkFlushInterval" yaml:"sinkFlushInterval"`
}

// DefaultConfig returns production-ready default configurations.
func DefaultConfig() Config {
	return Config{
		ChainID:           1,
		BlockChunkSize:    1000,
		MinBlockChunkSize: 1,
		MaxBlockChunkSize: 5000,
		Concurrency:       8,
		Confirmations:     12,
		PollInterval:      2 * time.Second,
		ChannelBufferSize: 20000,
		SinkBatchSize:     1000,
		SinkFlushInterval: 50 * time.Millisecond,
	}
}
