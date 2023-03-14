package cli

import (
	"flag"
)

type options struct {
	address string
	filter  bool
	listen  bool
	from    int
	to      int
}

func ParseArgs() *options {
	address := flag.String("addr", "", "Address of the ethereum transaction.")
	filter := flag.Bool("filter", true, "Filter and Index existing event logs on contract.")
	listen := flag.Bool("listen", true, "Listen for future events and index them for given contract.")
	from := flag.Int("from", 0, "Block range for querying, default value is genesis block.")
	to := flag.Int("to", 0, "Block range for querying, default value is latest block.")

	return &options{
		address: *address,
		filter:  *filter,
		listen:  *listen,
		from:    *from,
		to:      *to,
	}
}
