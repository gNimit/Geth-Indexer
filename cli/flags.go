// Package cli is used for parsing user input and reading config data.
package cli

import (
	"flag"
)

type FlagOptions struct {
	Address string
	From    int
	To      int
}

func ParseFlags() FlagOptions {
	address := flag.String("address", "", "Address of smart contract.")
	from := flag.Int("from", 0, "Block range for querying, default value is genesis block.")
	to := flag.Int("to", 0, "Block range for querying, default option makes application listen for future events.")

	flag.Parse()
	return FlagOptions{
		Address: *address,
		From:    *from,
		To:      *to,
	}
}
