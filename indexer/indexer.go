package indexer

import (
	"eventIndexer.com/subscriber"
	"fmt"
)

func Index(eventCh chan *subscriber.Event) {
	ed := <-eventCh
	fmt.Println(ed)
}
