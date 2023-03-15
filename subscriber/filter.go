// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
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

func filter(client *ethclient.Client, opts *cli.Options) []types.Log {

	if opts.FLAGS.From != 0 {
		from = big.NewInt(int64(opts.FLAGS.From))
	} else {
		from = nil
	}

	if opts.FLAGS.To != 0 {
		to = big.NewInt(int64(opts.FLAGS.To))
	} else {
		to = nil
	}

	query := ethereum.FilterQuery{
		FromBlock: from,
		ToBlock:   to,
		Addresses: []common.Address{
			common.HexToAddress(opts.FLAGS.Address),
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	return logs

}
