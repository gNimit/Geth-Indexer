package subscriber

import (
	"fmt"

	"eventIndexer.com/cli"
)

func Subscribe(events []string, opts *cli.Options) {
	client := newClient()
	defer client.Close()

	fmt.Printf("Subscribing to these events on contract... %v\n", events)
	if opts.FLAGS.Filter {
		logs := filter(client, opts)
		unpackFilterlogs(logs, events)
	}
	if opts.FLAGS.Listen {
		track(client, events, opts)
	}
}
