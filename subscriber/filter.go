// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var address = common.HexToAddress("0xdBFC942264f5CebF8C59f4065af2EFfB92D12475")

/*
	type queryOpts struct {
		FromBlock *big.Int
		ToBlock   *big.Int
		Address   *common.Address
	}
*/
func filter(client *ethclient.Client) {

	query := ethereum.FilterQuery{
		FromBlock: nil,
		ToBlock:   nil,
		Addresses: []common.Address{
			address,
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		log.Fatal(err)
	}

	unpackLogs(logs)

}
