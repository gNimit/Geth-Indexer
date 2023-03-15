package main

import (
	"flag"
	"log"
	"os"

	"eventIndexer.com/cli"
	"eventIndexer.com/subscriber"
)

func main() {
	exec()
}

func exec() {
	opts := cli.Run()
	events := flag.Args()
	if len(events) == 0 {
		log.Println("No events provided, please specify events as command line arguments. Exiting...")
		os.Exit(1)
	}

	subscriber.Subscribe(events, opts)
}
