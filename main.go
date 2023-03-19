package main

import (
	"eventIndexer.com/cli"
	"eventIndexer.com/indexer"
	"eventIndexer.com/subscriber"
	"flag"
	"log"
	"os"
)

func main() {
	os.Exit(exec())
}

func exec() int {
	opts := cli.Run()
	events := flag.Args()
	eventCh := make(chan *subscriber.Event, 1000)
	if len(events) == 0 {
		log.Println("No events provided, please specify smart contract events you want application to subscribe to. Exiting Application...")
		return 1
	}

	subscriber.Subscribe(events, eventCh, opts)
	db, err := indexer.Connect(opts.Database)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer func() int {
		if err := db.Close(); err != nil {
			log.Println(err)
			return 1
		}
		return 0
	}()

	go indexer.Index(<-eventCh)
	return 0
}
