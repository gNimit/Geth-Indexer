// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Generates a new ethclient.Client type with connection to given URL.
func newClient(infuraURL string) *ethclient.Client {

	client, err := ethclient.DialContext(context.Background(), infuraURL)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
