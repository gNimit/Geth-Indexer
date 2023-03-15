package subscriber

import (
	"fmt"
	"log"
	"strings"

	"eventIndexer.com/indexer"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	etherScanURL = "https://api-goerli.etherscan.io/api?module=contract&action=getabi&address=0xdBFC942264f5CebF8C59f4065af2EFfB92D12475&apikey=MXG89SYP3I7CBRWEJ5DJR1N2618YV4N7VC"
	contractABI  abi.ABI
)

func init() {
	contractABI = getContractABI()
}

func unpackFilterlogs(logs []types.Log, events []string) {
	for _, event := range events {
		for _, vLog := range logs {
			unpackLog(event, vLog.Data)
		}
	}
}

func unpackTrackLogs(vLog types.Log, events []string) {
	for _, event := range events {
		unpackLog(event, vLog.Data)
	}
}

func unpackLog(event string, data []byte) {
	var ifc = make(map[string]interface{})
	err := contractABI.UnpackIntoMap(ifc, event, data)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(ifc)
		publish(ifc, event)
	}
}

func getContractABI() abi.ABI {
	contractABI, err := abi.JSON(strings.NewReader(fetchABI(etherScanURL)))
	if err != nil {
		log.Println(err)
	}
	return contractABI
}

func publish(vLog map[string]interface{}, event string) {
	indexer.Index(event, vLog)
}
