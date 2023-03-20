package indexer

import (
	"database/sql"
	"eventIndexer.com/subscriber"
	"fmt"
	"strings"
)

func Index(eventCh chan *subscriber.Event, db *sql.DB) {

	for {
		ed := <-eventCh
		query := mapToQuery(strings.ToLower(ed.Name), ed)
		go upload(db, query)
	}
}

func mapToQuery(table string, params *subscriber.Event) string {
	columns := len(params.Data)
	fieldSlice := make([]string, 0, columns)
	fields := "name, blockNumber, blockHas, contract, "
	for field := range params.Data {
		fieldSlice = append(fieldSlice, field)
	}
	fields += strings.Join(fieldSlice, " ,")

	values := insertValues(fieldSlice, params)
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, fields, values)
}

func insertValues(keys []string, params *subscriber.Event) string {
	values := make([]string, 0, len(keys)+4)
	values = append(values, fmt.Sprintf("%v", params.Name), fmt.Sprintf("%v", params.BlockNumber),
		fmt.Sprintf("%v", params.BlockHash), fmt.Sprintf("%v", params.Contract))

	for _, k := range keys {
		values = append(values, fmt.Sprintf("%v", params.Data[k]))
	}
	return strings.Join(values, " ,")
}
