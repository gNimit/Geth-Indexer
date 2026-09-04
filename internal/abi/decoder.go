package abi

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gNimit/geth-indexer/pkg/types"
)

// EventDefinition pairs an ABI event with its precomputed indexed/non-indexed arguments.
type EventDefinition struct {
	Event      abi.Event
	ABI        abi.ABI
	Contract   common.Address
	Indexed    abi.Arguments
	NonIndexed abi.Arguments
}

// Registry manages contract ABIs and decodes raw logs into typed Event structs.
type Registry struct {
	mu sync.RWMutex

	// contractEvents maps (contractAddress + topic0) -> EventDefinition
	contractEvents map[common.Address]map[common.Hash]*EventDefinition

	// globalEvents maps topic0 -> EventDefinition (fallback for wildcard contracts)
	globalEvents map[common.Hash]*EventDefinition

	// allTopics holds all registered topic0 hashes for query filtering
	allTopics []common.Hash
}

// NewRegistry initializes an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		contractEvents: make(map[common.Address]map[common.Hash]*EventDefinition),
		globalEvents:   make(map[common.Hash]*EventDefinition),
	}
}

// RegisterContract parses an ABI (from string or reader) and indexes the requested events.
// If allowedEvents is empty, all events defined in the ABI are indexed.
func (r *Registry) RegisterContract(contract common.Address, abiJSON string, allowedEvents []string) error {
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return fmt.Errorf("abi: failed to parse ABI for %s: %w", contract.Hex(), err)
	}
	return r.RegisterParsedContract(contract, &parsedABI, allowedEvents)
}

// RegisterParsedContract registers a pre-parsed abi.ABI object.
func (r *Registry) RegisterParsedContract(contract common.Address, parsedABI *abi.ABI, allowedEvents []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filterMap := make(map[string]bool, len(allowedEvents))
	for _, e := range allowedEvents {
		filterMap[e] = true
	}

	if _, ok := r.contractEvents[contract]; !ok {
		r.contractEvents[contract] = make(map[common.Hash]*EventDefinition)
	}

	for _, ev := range parsedABI.Events {
		if len(allowedEvents) > 0 && !filterMap[ev.Name] {
			continue
		}

		var indexedArgs, nonIndexedArgs abi.Arguments
		for _, arg := range ev.Inputs {
			if arg.Indexed {
				indexedArgs = append(indexedArgs, arg)
			} else {
				nonIndexedArgs = append(nonIndexedArgs, arg)
			}
		}

		def := &EventDefinition{
			Event:      ev,
			ABI:        *parsedABI,
			Contract:   contract,
			Indexed:    indexedArgs,
			NonIndexed: nonIndexedArgs,
		}

		r.contractEvents[contract][ev.ID] = def
		r.globalEvents[ev.ID] = def

		// Check if topic is already in allTopics
		exists := false
		for _, t := range r.allTopics {
			if t == ev.ID {
				exists = true
				break
			}
		}
		if !exists {
			r.allTopics = append(r.allTopics, ev.ID)
		}
	}

	return nil
}

// Topics returns all registered event signature topic hashes.
func (r *Registry) Topics() []common.Hash {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]common.Hash, len(r.allTopics))
	copy(out, r.allTopics)
	return out
}

// Decode transforms a raw Ethereum log into a decoded Event.
// If the log does not match registered events, it returns (nil, nil).
func (r *Registry) Decode(log ethTypes.Log, chainID uint64) (*types.Event, error) {
	if len(log.Topics) == 0 {
		return nil, nil // Anonymous logs without topics cannot be matched by topic0
	}

	topic0 := log.Topics[0]

	r.mu.RLock()
	var def *EventDefinition
	if eventsForContract, ok := r.contractEvents[log.Address]; ok {
		def = eventsForContract[topic0]
	}
	if def == nil {
		def = r.globalEvents[topic0]
	}
	r.mu.RUnlock()

	if def == nil {
		return nil, nil // Ignored event
	}

	dataMap := make(map[string]interface{})

	// 1. Unpack non-indexed fields from log.Data
	if len(log.Data) > 0 && len(def.NonIndexed) > 0 {
		if err := def.NonIndexed.UnpackIntoMap(dataMap, log.Data); err != nil {
			if err2 := def.ABI.UnpackIntoMap(dataMap, def.Event.Name, log.Data); err2 != nil {
				return nil, fmt.Errorf("abi: failed to unpack non-indexed log data for %s: %w", def.Event.Name, err)
			}
		}
	}

	// 2. Unpack indexed fields from Topics[1:]
	if len(log.Topics) > 1 && len(def.Indexed) > 0 {
		if err := abi.ParseTopicsIntoMap(dataMap, def.Indexed, log.Topics[1:]); err != nil {
			return nil, fmt.Errorf("abi: failed to parse indexed topics for %s: %w", def.Event.Name, err)
		}
	}

	topicsCopy := make([]common.Hash, len(log.Topics))
	copy(topicsCopy, log.Topics)

	dataCopy := make([]byte, len(log.Data))
	copy(dataCopy, log.Data)

	event := &types.Event{
		ChainID:         chainID,
		ContractAddress: log.Address,
		EventName:       def.Event.Name,
		EventSignature:  topic0,
		BlockNumber:     log.BlockNumber,
		BlockHash:       log.BlockHash,
		TxHash:          log.TxHash,
		TxIndex:         log.TxIndex,
		LogIndex:        log.Index,
		Removed:         log.Removed,
		Data:            dataMap,
		RawData:         dataCopy,
		RawTopics:       topicsCopy,
		Timestamp:       time.Now().UTC(),
	}

	return event, nil
}
