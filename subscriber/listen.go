// Package subscriber Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"context"
	"eventIndexer.com/cli"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func listen(client *ethclient.Client, opts *cli.Config) ethereum.Subscription {
	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(opts.Query.Address)},
	}

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Printf("Failed to subscriber to logs, %v\n", err)
	}
	return sub
}
