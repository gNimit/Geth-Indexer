// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func track(client *ethclient.Client) {
	logs := make(chan types.Log)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Printf("Failed to subscriber to logs, %v\n", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription Error occured, %v\n", err)
		case vLog := <-logs:
			fmt.Println(vLog)
		}
	}
}
