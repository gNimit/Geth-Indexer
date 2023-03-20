// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"context"
	"log"
	"math/big"

	"eventIndexer.com/cli"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	from *big.Int
	to   *big.Int
)

func filter(client *ethclient.Client, opts *cli.Config) []types.Log {

	if opts.Query.From != 0 {
		from = big.NewInt(int64(opts.Query.From))
	} else {
		from = nil
	}

	if opts.Query.To != 0 {
		to = big.NewInt(int64(opts.Query.To))
	} else {
		to = nil
	}

	query := ethereum.FilterQuery{
		FromBlock: from,
		ToBlock:   to,
		Addresses: []common.Address{
			common.HexToAddress(opts.Query.Address),
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	return logs
}
