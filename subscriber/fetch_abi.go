// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var abiMap = make(map[string]string)

func fetchABI(etherScanURL string) string {
	resp, err := http.Get(etherScanURL)
	if err != nil {
		log.Printf("Failed to fetch ABI from etherscan : %v\n", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read from http response on ABI : %v\n", err)
	}

	unmarshalToMap(data, abiMap)

	return abiMap["result"]
}

func unmarshalToMap(data []byte, abiMap map[string]string) {
	err := json.Unmarshal(data, &abiMap)
	if err != nil {
		log.Printf("Failed to unmarshal data into map : %v\n", err)
	}
}
