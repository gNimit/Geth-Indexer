// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	address      = common.HexToAddress("0xdBFC942264f5CebF8C59f4065af2EFfB92D12475")
	etherScanAPI = "https://api-goerli.etherscan.io/api?module=contract&action=getabi&address=0xdBFC942264f5CebF8C59f4065af2EFfB92D12475&apikey=MXG89SYP3I7CBRWEJ5DJR1N2618YV4N7VC"
)

type queryOpts struct {
	FromBlock *big.Int
	ToBlock   *big.Int
	Address   *common.Address
}

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

	contractABI, err := abi.JSON(strings.NewReader(fetchABI(etherScanAPI)))
	if err != nil {
		log.Println(err)
	}

	var ifc = make(map[string]interface{})
	for _, vLog := range logs {
		err := contractABI.UnpackIntoMap(ifc, "LotJoined", vLog.Data)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(ifc)
	}

}
