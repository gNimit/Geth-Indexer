// Package to deal with all event subscription functionalities on a ethereum contract.
package subscriber

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ctx       = context.Background()
	infuraURL = "wss://goerli.infura.io/ws/v3/818b966b2910464081ce9f7ffabe4092"
)

// Generates a new ethclient.Client type with connection to given URL.
func newClient() *ethclient.Client {

	client, err := ethclient.DialContext(ctx, infuraURL)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
