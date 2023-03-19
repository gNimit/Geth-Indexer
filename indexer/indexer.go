package indexer

import (
	"eventIndexer.com/subscriber"
	"fmt"
)

type eventData struct {
	event string
	log   map[string]interface{}
}

func Index(ed *subscriber.Event) {
	fmt.Println(ed)
}
