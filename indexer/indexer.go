package indexer

import "eventIndexer.com/cli"

type eventData struct {
	event string
	log   map[string]interface{}
}

func Index(_event string, _vLog map[string]interface{}, dbOps cli.DbOptions) {

	db := connect(dbOps)
	defer db.Close()

	ed := eventData{
		event: _event,
		log:   _vLog,
	}
	upload(db, ed)
}
