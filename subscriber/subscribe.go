package subscriber

import (
	"eventIndexer.com/cli"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"strings"
)

const chanBufferSize = 1000

type Contract struct {
	Address common.Address
	ABI     abi.ABI
	events  map[common.Hash]string
}

type Event struct {
	Name        string
	BlockNumber uint64
	BlockHash   common.Hash
	Contract    common.Address
	Data        map[string]interface{}
}

func Subscribe(events []string, eventCh chan *Event, opts *cli.Config, quit chan struct{}) {
	client := newClient(opts.API.EthNodeURL)
	defer client.Close()

	fmt.Printf("Subscribing to these events on contract %s ... %s\n", opts.Query.Address, strings.Join(events, " "))
	c := &Contract{
		Address: common.HexToAddress(opts.Query.Address),
		ABI:     fetchABI(opts.API.EtherScanAPI),
		events:  make(map[common.Hash]string),
	}

	for _, e := range c.ABI.Events {
		c.events[e.ID] = e.Name
	}

	logCh := make(chan types.Log, chanBufferSize)

	go func() {
		for _, l := range filter(client, opts) {
			logCh <- l
		}
	}()

	sub := listen(client, opts)

	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
		case l := <-logCh:
			if edata := parseEvents(events, l, c); edata != nil {
				eventCh <- edata
			}
		case <-quit:
			break
		}
	}

}

// Establishes connection to ethereum node and returns ethclient.Client type.
func newClient(EthNodeURL string) *ethclient.Client {

	client, err := ethclient.Dial(EthNodeURL)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func parseEvents(events []string, l types.Log, c *Contract) *Event {
	name, ok := c.events[l.Topics[0]]
	if !ok {
		return nil
	}
	event := ""

	for _, e := range events {
		if e == name {
			event = name
			break
		}
	}
	if event == "" {
		return nil
	}

	data, err := unpackLog(event, l.Data, c.ABI)
	if err != nil || data == nil {
		return nil
	}

	e := &Event{
		Name:        name,
		BlockNumber: l.BlockNumber,
		BlockHash:   l.BlockHash,
		Contract:    c.Address,
		Data:        data,
	}
	return e

}

func unpackLog(event string, data []byte, contractABI abi.ABI) (map[string]interface{}, error) {
	logMap := make(map[string]interface{})
	err := contractABI.UnpackIntoMap(logMap, event, data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return logMap, nil
}
