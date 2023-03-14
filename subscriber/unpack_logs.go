package subscriber

import (
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

var etherScanURL = "https://api-goerli.etherscan.io/api?module=contract&action=getabi&address=0xdBFC942264f5CebF8C59f4065af2EFfB92D12475&apikey=MXG89SYP3I7CBRWEJ5DJR1N2618YV4N7VC"

func unpackLogs(logs []types.Log) {

	contractABI, err := abi.JSON(strings.NewReader(fetchABI(etherScanURL)))
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
