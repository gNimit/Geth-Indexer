package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gNimit/geth-indexer/pkg/checkpoint"
	"github.com/gNimit/geth-indexer/pkg/indexer"
	"github.com/gNimit/geth-indexer/pkg/metrics"
	"github.com/gNimit/geth-indexer/pkg/rpc"
	"github.com/gNimit/geth-indexer/pkg/sink"
)

func main() {
	var (
		rpcURLs        = flag.String("rpc", "http://127.0.0.1:8545", "Comma-separated RPC endpoint URLs")
		chainID        = flag.Uint64("chain-id", 1, "Ethereum Chain ID")
		contractAddr   = flag.String("contract", "", "Target smart contract address")
		abiPath        = flag.String("abi", "", "Path to contract ABI JSON file")
		eventNames     = flag.String("events", "", "Comma-separated list of event names to index (empty = all)")
		startBlock     = flag.Uint64("from", 0, "Starting block number")
		endBlock       = flag.Uint64("to", 0, "Ending block number (0 = continuous live tailing)")
		concurrency    = flag.Int("concurrency", 8, "Number of concurrent block fetching workers")
		chunkSize      = flag.Uint64("chunk-size", 1000, "Initial block chunk size")
		confirmations  = flag.Uint64("confirmations", 12, "Reorg safety confirmation blocks")
		pollInterval   = flag.Duration("poll-interval", 2*time.Second, "Polling interval for live streaming")
		sinkType       = flag.String("sink", "stdout", "Sink type: stdout, postgres, webhook")
		pgHost         = flag.String("pg-host", "localhost", "PostgreSQL host")
		pgPort         = flag.Int("pg-port", 5432, "PostgreSQL port")
		pgUser         = flag.String("pg-user", "postgres", "PostgreSQL user")
		pgPass         = flag.String("pg-pass", "postgres", "PostgreSQL password")
		pgDB           = flag.String("pg-db", "geth_indexer", "PostgreSQL database name")
		checkpointFile = flag.String("checkpoint-file", "indexer_checkpoint.json", "Path to JSON checkpoint file")
		metricsAddr    = flag.String("metrics-addr", ":9090", "Prometheus metrics HTTP server address")
	)
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	slog.Info("Starting geth-indexer daemon",
		"chainID", *chainID,
		"fromBlock", *startBlock,
		"toBlock", *endBlock,
		"concurrency", *concurrency,
		"chunkSize", *chunkSize,
	)

	// 1. Setup Metrics Collector & HTTP server
	metricCollector := metrics.NewPrometheusCollector("geth_indexer")
	if *metricsAddr != "" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", metricCollector.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		metricsServer := &http.Server{
			Addr:    *metricsAddr,
			Handler: mux,
		}
		go func() {
			slog.Info("Prometheus metrics server listening", "addr", *metricsAddr)
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("Metrics server failed", "error", err)
			}
		}()
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_ = metricsServer.Shutdown(ctx)
		}()
	}

	// 2. Setup RPC Pool
	urls := strings.Split(*rpcURLs, ",")
	var epConfigs []rpc.EndpointConfig
	for i, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		epConfigs = append(epConfigs, rpc.EndpointConfig{
			URL:      u,
			Provider: fmt.Sprintf("rpc-node-%d", i+1),
			Weight:   1,
			Timeout:  10 * time.Second,
		})
	}
	if len(epConfigs) == 0 {
		slog.Error("No valid RPC endpoints provided")
		os.Exit(1)
	}

	pool, err := rpc.NewPool(rpc.PoolConfig{
		ChainID:   *chainID,
		Endpoints: epConfigs,
		Balancer:  rpc.NewLeastInFlightBalancer(),
		Metrics:   metricCollector,
	}, nil)
	if err != nil {
		slog.Error("Failed to initialize RPC pool", "error", err)
		os.Exit(1)
	}

	// 3. Setup Checkpoint Store
	store, err := checkpoint.NewFileStore(*checkpointFile)
	if err != nil {
		slog.Error("Failed to initialize file checkpoint store", "error", err)
		os.Exit(1)
	}

	// 4. Setup Contract Configurations
	var contracts []indexer.ContractConfig
	if *contractAddr != "" {
		targetAddr := common.HexToAddress(*contractAddr)
		var abiJSON string
		if *abiPath != "" {
			data, err := os.ReadFile(*abiPath)
			if err != nil {
				slog.Error("Failed to read ABI file", "path", *abiPath, "error", err)
				os.Exit(1)
			}
			abiJSON = string(data)
		}

		var events []string
		if *eventNames != "" {
			for _, e := range strings.Split(*eventNames, ",") {
				if trimmed := strings.TrimSpace(e); trimmed != "" {
					events = append(events, trimmed)
				}
			}
		}

		contracts = append(contracts, indexer.ContractConfig{
			Address: targetAddr,
			ABIJSON: abiJSON,
			Events:  events,
		})
	}

	// 5. Build Indexer Config & Options
	cfg := indexer.Config{
		ChainID:           *chainID,
		Contracts:         contracts,
		StartBlock:        *startBlock,
		EndBlock:          *endBlock,
		BlockChunkSize:    *chunkSize,
		Concurrency:       *concurrency,
		Confirmations:     *confirmations,
		PollInterval:      *pollInterval,
		SinkBatchSize:     1000,
		SinkFlushInterval: 50 * time.Millisecond,
	}

	opts := []indexer.Option{
		indexer.WithRPCPool(pool),
		indexer.WithCheckpointStore(store),
		indexer.WithMetrics(metricCollector),
	}

	// Setup Selected Sink
	switch *sinkType {
	case "postgres":
		opts = append(opts, indexer.WithPostgresSink(sink.PostgresConfig{
			Host:     *pgHost,
			Port:     *pgPort,
			User:     *pgUser,
			Password: *pgPass,
			Database: *pgDB,
		}))
		slog.Info("PostgreSQL batch sink enabled", "host", *pgHost, "db", *pgDB)
	case "stdout":
		opts = append(opts, indexer.WithSink(sink.NewStdoutSink(false)))
	default:
		opts = append(opts, indexer.WithSink(sink.NewStdoutSink(false)))
	}

	idx, err := indexer.New(cfg, opts...)
	if err != nil {
		slog.Error("Failed to initialize indexer", "error", err)
		os.Exit(1)
	}

	// 6. Handle Signals for Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		slog.Info("Received termination signal, shutting down gracefully...", "signal", sig)
		cancel()
		if err := idx.Stop(); err != nil {
			slog.Error("Error during indexer shutdown", "error", err)
		}
	}()

	// 7. Start Indexing
	eventCh, err := idx.Start(ctx)
	if err != nil {
		slog.Error("Failed to start indexer", "error", err)
		os.Exit(1)
	}

	slog.Info("Indexer started successfully, listening for events...")

	// Drain events in background
	go func() {
		for ev := range eventCh {
			_ = ev
		}
	}()

	// Block until completion or context cancellation
	if err := idx.Wait(); err != nil && err != context.Canceled {
		slog.Error("Indexer terminated with error", "error", err)
		os.Exit(1)
	}

	slog.Info("geth-indexer exited cleanly.")
}
