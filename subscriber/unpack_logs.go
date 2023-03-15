package subscriber

import (
	"fmt"
	"log"
	"strings"

	"eventIndexer.com/cli"
	"eventIndexer.com/indexer"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

func unpackFilterlogs(logs []types.Log, events []string, opts *cli.Options) {
	etherScanURL := opts.API.ETHERSCAN_URL
	contractABI := getContractABI(etherScanURL)
	for _, event := range events {
		for _, vLog := range logs {
			unpackLog(event, vLog.Data, contractABI, opts.DB)
		}
	}
}

func unpackTrackLogs(vLog types.Log, events []string, opts *cli.Options) {
	etherScanURL := opts.API.ETHERSCAN_URL
	contractABI := getContractABI(etherScanURL)
	for _, event := range events {
		unpackLog(event, vLog.Data, contractABI, opts.DB)
	}
}

func unpackLog(event string, data []byte, contractABI abi.ABI, dbOps cli.DbOptions) {
	var ifc = make(map[string]interface{})
	err := contractABI.UnpackIntoMap(ifc, event, data)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(ifc)
		publish(ifc, event, dbOps)
	}
}

func getContractABI(etherScanURL string) abi.ABI {
	contractABI, err := abi.JSON(strings.NewReader(fetchABI(etherScanURL)))
	if err != nil {
		log.Println(err)
	}
	return contractABI
}

func publish(vLog map[string]interface{}, event string, dbOps cli.DbOptions) {
	indexer.Index(event, vLog, dbOps)
}
