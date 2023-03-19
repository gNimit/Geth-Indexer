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
}

type Event struct {
	Name        string
	BlockNumber uint64
	BlockHash   common.Hash
	Contract    common.Address
	Data        map[string]interface{}
}

// Establishes connection to ethereum node and returns ethclient.Client type.
func newClient(EthNodeURL string) *ethclient.Client {

	client, err := ethclient.Dial(EthNodeURL)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func Subscribe(events []string, eventCh chan<- *Event, opts *cli.Config) {
	client := newClient(opts.API.EthNodeURL)
	defer client.Close()

	fmt.Printf("Subscribing to these events on contract %s ... %s\n", opts.Flag.Address, strings.Join(events, " "))

	c := &Contract{
		Address: common.HexToAddress(opts.Flag.Address),
		ABI:     fetchABI(opts.API.EtherScanURL),
	}
	for _, event := range events {

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
				edata := parseEvents(event, l, c)
				eventCh <- edata
			}
		}
	}
}

func parseEvents(event string, l types.Log, c *Contract) *Event {

	data, err := unpackLog(event, l.Data, c.ABI)
	if err != nil || data == nil {
		return nil
	}

	e := &Event{
		Name:        event,
		BlockNumber: l.BlockNumber,
		BlockHash:   l.BlockHash,
		Contract:    c.Address,
		Data:        data,
	}
	//fmt.Println(e)
	return e
}
