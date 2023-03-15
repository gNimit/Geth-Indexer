package indexer

type eventData struct {
	event string
	log   map[string]interface{}
}

func Index(_event string, _vLog map[string]interface{}) {

	db := connect()
	defer db.Close()

	ed := eventData{
		event: _event,
		log:   _vLog,
	}
	upload(db, ed)
}
