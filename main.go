package main

import (
	"bufio"
	"eventIndexer.com/cli"
	"eventIndexer.com/indexer"
	"eventIndexer.com/subscriber"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const chanBufferSize = 1000

func main() {
	os.Exit(exec())
}

func exec() int {
	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(2)

	opts := cli.Run()
	events := flag.Args()
	if len(events) == 0 {
		log.Println("No events provided, please specify smart contract events you want application to subscribe to. Exiting Application...")
		return 1
	}

	eventCh := make(chan *subscriber.Event, chanBufferSize)
	quitCh := make(chan bool)

	go stopSignal(quitCh)
	go subscriber.Subscribe(events, eventCh, opts, quitCh)

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
	go indexer.Index(eventCh, db, quitCh)

	wg.Wait()
	return 0
}

func stopSignal(quitCh chan bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input stop to exit from application")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input, try again or press ctrl+c to exit from application", err)
	}
	if strings.TrimSpace(text) == "stop" {
		quitCh <- true
	}
}
