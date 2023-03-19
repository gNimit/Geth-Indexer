package subscriber

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"log"
)

func unpackLog(event string, data []byte, contractABI abi.ABI) (map[string]interface{}, error) {
	var logMap = make(map[string]interface{})
	err := contractABI.UnpackIntoMap(logMap, event, data)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	return logMap, nil
}
