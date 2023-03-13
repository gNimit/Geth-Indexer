// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"context"
	"fmt"
	"log"

	multiLot "eventIndexer.com/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

func listenEvents(client *ethclient.Client) {
	tc, err := multiLot.NewMultiLot(address, client)
	if err != nil {
		log.Printf("Failed to establish connection with transaction at address %v, Error : %v\n", address, err)
	}

	watchOpts := &bind.WatchOpts{Context: context.Background(), Start: nil}
	listenLotCreated(tc, watchOpts)
	listenLotJoined(tc, watchOpts)
	listenLotResolved(tc, watchOpts)
}

func listenLotCreated(tc *multiLot.MultiLot, watchOpts *bind.WatchOpts) {
	channel := make(chan *multiLot.MultiLotLotCreated)
	go func() {
		sub, err := tc.WatchLotCreated(watchOpts, channel)
		if err != nil {
			log.Printf("Failed to watch LotCreated Event %v", err)
			defer sub.Unsubscribe()
		}
	}()

	event := <-channel
	fmt.Println(event.Creator.Hex())
}

func listenLotJoined(tc *multiLot.MultiLot, watchOpts *bind.WatchOpts) {
	channel := make(chan *multiLot.MultiLotLotJoined)
	go func() {
		sub, err := tc.WatchLotJoined(watchOpts, channel)
		if err != nil {
			log.Printf("Failed to watch LotCreated Event %v", err)
			defer sub.Unsubscribe()
		}
	}()

	event := <-channel
	fmt.Println(event.User.Hex())
}

func listenLotResolved(tc *multiLot.MultiLot, watchOpts *bind.WatchOpts) {
	channel := make(chan *multiLot.MultiLotLotResolved)
	go func() {
		sub, err := tc.WatchLotResolved(watchOpts, channel)
		if err != nil {
			log.Printf("Failed to watch LotCreated Event %v", err)
			defer sub.Unsubscribe()
		}
	}()

	event := <-channel
	fmt.Println(event.Raw.Address.Hex())
}
