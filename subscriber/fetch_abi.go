// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io"
	"log"
	"net/http"
	"strings"
)

func fetchABI(etherScanURL string) abi.ABI {
	resp, err := http.Get(etherScanURL)
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
