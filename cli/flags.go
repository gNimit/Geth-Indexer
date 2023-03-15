package cli

import (
	"flag"
)

type FlagOptions struct {
	Address string
	Filter  bool
	Listen  bool
	From    int
	To      int
}

func ParseArgs() FlagOptions {
	address := flag.String("address", "", "Address of the ethereum transaction.")
	filter := flag.Bool("filter", true, "Filter and Index existing event logs on contract.")
	listen := flag.Bool("listen", true, "Listen for future events and index them for given contract.")
	from := flag.Int("from", 0, "Block range for querying, default value is genesis block.")
	to := flag.Int("to", 0, "Block range for querying, default value is latest block.")

	flag.Parse()
	return FlagOptions{
		Address: *address,
		Filter:  *filter,
		Listen:  *listen,
		From:    *from,
		To:      *to,
	}
}
