package main

import (
	"fmt"
	"log"
	"os"

	"eventIndexer.com/cli"
	"eventIndexer.com/subscriber"
)

func main() {

	if len(os.Args) == 0 {
		log.Println("No events provided, please specify events as command line arguments. Exiting...")
		os.Exit(1)
	}
	run()
}

func run() {
	events := os.Args[1:]
	opts := cli.ParseArgs()
	fmt.Println(events, opts)
	subscriber.Subscribe()
}
