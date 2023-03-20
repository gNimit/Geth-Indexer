// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io"
	"log"
	"net/http"
	"strings"
)

const etherScanURL = "https://api-goerli.etherscan.io/api?module=contract&action=getabi&address=0xdBFC942264f5CebF8C59f4065af2EFfB92D12475&apikey=%s"

func fetchABI(etherScanAPI string) abi.ABI {
	url := fmt.Sprintf(etherScanURL, etherScanAPI)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch ABI from etherscan : %v\n", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read from http response on ABI : %v\n", err)
	}

	a, err := abi.JSON(strings.NewReader(unmarshalToMap(data)))
	if err != nil {
		log.Println(err)
	}
	return a
}

func unmarshalToMap(data []byte) string {
	abiMap := make(map[string]string)
	err := json.Unmarshal(data, &abiMap)
	if err != nil {
		log.Printf("Failed to unmarshal data into map : %v\n", err)
	}
	return abiMap["result"]
}
