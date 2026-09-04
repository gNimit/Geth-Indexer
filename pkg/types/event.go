package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Event represents a decoded smart contract event log emitted on-chain.
type Event struct {
	ChainID         uint64                 `json:"chainId"`
	ContractAddress common.Address         `json:"contractAddress"`
	EventName       string                 `json:"eventName"`
	EventSignature  common.Hash            `json:"eventSignature"` // Keccak-256 hash of event signature (Topic[0])
	BlockNumber     uint64                 `json:"blockNumber"`
	BlockHash       common.Hash            `json:"blockHash"`
	TxHash          common.Hash            `json:"txHash"`
	TxIndex         uint                   `json:"txIndex"`
	LogIndex        uint                   `json:"logIndex"`
	Removed         bool                   `json:"removed"`
	Data            map[string]interface{} `json:"data"`
	RawData         []byte                 `json:"rawData,omitempty"`
	RawTopics       []common.Hash          `json:"rawTopics,omitempty"`
	Timestamp       time.Time              `json:"timestamp"`
}

// ID returns a globally unique composite identifier for this event log.
func (e *Event) ID() string {
	return fmt.Sprintf("%d:%s:%d", e.ChainID, e.TxHash.Hex(), e.LogIndex)
}

// JSON serializes the event into JSON bytes.
func (e *Event) JSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) String() string {
	return fmt.Sprintf("[%d] Block %d | %s | %s | Tx %s | Log %d",
		e.ChainID, e.BlockNumber, e.ContractAddress.Hex(), e.EventName, e.TxHash.Hex(), e.LogIndex)
}
