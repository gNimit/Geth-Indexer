package indexer

import (
	"container/heap"
	"context"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/types"
)

type blockChunk struct {
	FromBlock uint64
	ToBlock   uint64
	IsGap     bool
	GapRange  checkpoint.BlockRange
}

type chunkResult struct {
	Chunk  blockChunk
	Events []*types.Event
	Err    error
}

// resultHeap implements a min-heap ordered by Chunk.FromBlock for deterministic monotonic event ordering.
type resultHeap []*chunkResult

func (h resultHeap) Len() int           { return len(h) }
func (h resultHeap) Less(i, j int) bool { return h[i].Chunk.FromBlock < h[j].Chunk.FromBlock }
func (h resultHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *resultHeap) Push(x interface{}) {
	*h = append(*h, x.(*chunkResult))
}
func (h *resultHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (idx *Indexer) runPipeline(ctx context.Context) error {
	addresses := idx.ContractAddresses()
	topics := idx.registry.Topics()

	var topicFilter [][]common.Hash
	if len(topics) > 0 {
		topicFilter = [][]common.Hash{topics}
	}

	// 1. Determine Starting Block and check for historical gaps
	startBlock := idx.config.StartBlock
	for _, addr := range addresses {
		lastBlock, err := idx.store.GetLastIndexedBlock(ctx, idx.config.ChainID, addr)
		if err == nil && lastBlock > 0 {
			if lastBlock+1 > startBlock {
				startBlock = lastBlock + 1
			}
		}
	}

	// 2. Initialize worker channels
	taskCh := make(chan blockChunk, idx.config.Concurrency*2)
	resultCh := make(chan *chunkResult, idx.config.Concurrency*2)

	var workerWg sync.WaitGroup
	for w := 0; w < idx.config.Concurrency; w++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			idx.workerLoop(ctx, taskCh, resultCh, addresses, topicFilter)
		}()
	}

	// 3. Coordinator goroutine: receives results, orders monotonically, flushes batches to sinks
	coordinatorDone := make(chan error, 1)
	go func() {
		coordinatorDone <- idx.coordinatorLoop(ctx, resultCh, startBlock)
	}()

	// 4. Feeder: generates chunks (both pending gaps and continuous blocks)
	feederErr := idx.feederLoop(ctx, taskCh, startBlock, addresses)
	close(taskCh)
	workerWg.Wait()
	close(resultCh)

	coordErr := <-coordinatorDone
	if feederErr != nil && feederErr != context.Canceled {
		return feederErr
	}
	return coordErr
}

func (idx *Indexer) feederLoop(ctx context.Context, taskCh chan<- blockChunk, startBlock uint64, addresses []common.Address) error {
	// First: feed any pending unindexed gaps from checkpoint store
	for _, addr := range addresses {
		gaps, err := idx.store.GetUnindexedRanges(ctx, idx.config.ChainID, addr)
		if err == nil {
			for _, gap := range gaps {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case taskCh <- blockChunk{
					FromBlock: gap.FromBlock,
					ToBlock:   gap.ToBlock,
					IsGap:     true,
					GapRange:  gap,
				}:
				}
			}
		}
	}

	currentBlock := startBlock
	chunkSize := idx.config.BlockChunkSize

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		latestBlock, err := idx.pool.BlockNumber(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(500 * time.Millisecond):
				continue
			}
		}

		var safeHead uint64
		if latestBlock > idx.config.Confirmations {
			safeHead = latestBlock - idx.config.Confirmations
		} else {
			safeHead = 0
		}

		// Check if we have caught up to the safe head
		if currentBlock > safeHead {
			if idx.config.EndBlock > 0 && currentBlock > idx.config.EndBlock {
				return nil // Range indexing finished
			}
			// Live streaming mode: wait for new blocks
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(idx.config.PollInterval):
				continue
			}
		}

		toBlock := currentBlock + chunkSize - 1
		if toBlock > safeHead {
			toBlock = safeHead
		}
		if idx.config.EndBlock > 0 && toBlock > idx.config.EndBlock {
			toBlock = idx.config.EndBlock
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case taskCh <- blockChunk{FromBlock: currentBlock, ToBlock: toBlock}:
		}

		if idx.config.EndBlock > 0 && toBlock >= idx.config.EndBlock {
			return nil
		}
		currentBlock = toBlock + 1
	}
}

func (idx *Indexer) workerLoop(
	ctx context.Context,
	taskCh <-chan blockChunk,
	resultCh chan<- *chunkResult,
	addresses []common.Address,
	topicFilter [][]common.Hash,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case chunk, ok := <-taskCh:
			if !ok {
				return
			}

			logs, err := idx.fetchLogsAdaptive(ctx, chunk.FromBlock, chunk.ToBlock, addresses, topicFilter)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case resultCh <- &chunkResult{Chunk: chunk, Err: err}:
				}
				continue
			}

			events := make([]*types.Event, 0, len(logs))
			for _, l := range logs {
				ev, err := idx.registry.Decode(l, idx.config.ChainID)
				if err != nil {
					continue
				}
				if ev != nil {
					events = append(events, ev)
				}
			}

			select {
			case <-ctx.Done():
				return
			case resultCh <- &chunkResult{Chunk: chunk, Events: events}:
			}
		}
	}
}

// fetchLogsAdaptive executes eth_getLogs with automatic recursive range bisection on RPC limits.
func (idx *Indexer) fetchLogsAdaptive(
	ctx context.Context,
	from uint64,
	to uint64,
	addresses []common.Address,
	topics [][]common.Hash,
) ([]ethTypes.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: addresses,
		Topics:    topics,
	}

	logs, err := idx.pool.FilterLogs(ctx, query)
	if err == nil {
		return logs, nil
	}

	if isBlockRangeLimitError(err) && from < to {
		mid := from + (to-from)/2
		left, err1 := idx.fetchLogsAdaptive(ctx, from, mid, addresses, topics)
		if err1 != nil {
			return nil, err1
		}
		right, err2 := idx.fetchLogsAdaptive(ctx, mid+1, to, addresses, topics)
		if err2 != nil {
			return nil, err2
		}
		return append(left, right...), nil
	}

	return nil, err
}

func (idx *Indexer) coordinatorLoop(ctx context.Context, resultCh <-chan *chunkResult, expectedBlock uint64) error {
	var results resultHeap
	heap.Init(&results)

	currentExpected := expectedBlock
	eventBuffer := make([]*types.Event, 0, idx.config.SinkBatchSize)
	flushTicker := time.NewTicker(idx.config.SinkFlushInterval)
	defer flushTicker.Stop()

	flushEvents := func(fCtx context.Context) error {
		if len(eventBuffer) == 0 {
			return nil
		}
		start := time.Now()
		for _, s := range idx.sinks {
			if err := s.Send(fCtx, eventBuffer); err != nil {
				idx.metrics.IncSinkEvents(s.Name(), len(eventBuffer), "error")
				continue
			}
			idx.metrics.IncSinkEvents(s.Name(), len(eventBuffer), "success")
			idx.metrics.ObserveBatchProcessing(s.Name(), time.Since(start))
		}
		idx.metrics.IncEventsProcessed(len(eventBuffer))
		eventBuffer = eventBuffer[:0]
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			_ = flushEvents(shutdownCtx)
			shutdownCancel()
			return ctx.Err()

		case <-flushTicker.C:
			if err := flushEvents(ctx); err != nil {
				return err
			}

		case res, ok := <-resultCh:
			if !ok {
				// Result channel closed; flush remaining items in heap and buffer
				for results.Len() > 0 {
					item := heap.Pop(&results).(*chunkResult)
					eventBuffer = append(eventBuffer, item.Events...)
					_ = idx.commitProgress(ctx, item)
				}
				flushCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				return flushEvents(flushCtx)
			}

			if res.Err != nil {
				// Record unindexed range for future backfilling
				for _, c := range idx.config.Contracts {
					_ = idx.store.AddUnindexedRange(ctx, idx.config.ChainID, c.Address, checkpoint.BlockRange{
						FromBlock: res.Chunk.FromBlock,
						ToBlock:   res.Chunk.ToBlock,
					})
				}
				continue
			}

			// Sort events in this chunk deterministically by (BlockNumber, TxIndex, LogIndex)
			sort.Slice(res.Events, func(i, j int) bool {
				if res.Events[i].BlockNumber != res.Events[j].BlockNumber {
					return res.Events[i].BlockNumber < res.Events[j].BlockNumber
				}
				if res.Events[i].TxIndex != res.Events[j].TxIndex {
					return res.Events[i].TxIndex < res.Events[j].TxIndex
				}
				return res.Events[i].LogIndex < res.Events[j].LogIndex
			})

			heap.Push(&results, res)

			// Drain contiguous chunks from the min-heap
			for results.Len() > 0 {
				top := results[0]
				if top.Chunk.IsGap || top.Chunk.FromBlock <= currentExpected {
					item := heap.Pop(&results).(*chunkResult)
					eventBuffer = append(eventBuffer, item.Events...)
					if item.Chunk.ToBlock >= currentExpected {
						currentExpected = item.Chunk.ToBlock + 1
					}
					if err := idx.commitProgress(ctx, item); err != nil {
						return err
					}

					if len(eventBuffer) >= idx.config.SinkBatchSize {
						if err := flushEvents(ctx); err != nil {
							return err
						}
					}
				} else {
					break
				}
			}
		}
	}
}

func (idx *Indexer) commitProgress(ctx context.Context, item *chunkResult) error {
	for _, c := range idx.config.Contracts {
		if item.Chunk.IsGap {
			_ = idx.store.MarkRangeCompleted(ctx, idx.config.ChainID, c.Address, item.Chunk.GapRange)
		} else {
			_ = idx.store.SaveLastIndexedBlock(ctx, idx.config.ChainID, c.Address, item.Chunk.ToBlock)
		}
	}
	idx.metrics.SetCurrentBlock(idx.config.ChainID, item.Chunk.ToBlock)
	return nil
}

func isBlockRangeLimitError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "10000") ||
		strings.Contains(msg, "more than") ||
		strings.Contains(msg, "range") ||
		strings.Contains(msg, "too large") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "limit exceeded") ||
		strings.Contains(msg, "maximum block range")
}
