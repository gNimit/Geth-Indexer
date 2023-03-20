package indexer

import (
	"database/sql"
	"log"
)

func upload(db *sql.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		log.Println("Failed to execute ", query, err)
	}
}

func indexingCheck(db *sql.DB, relation string, column string) {
	index := relation + "_" + column + "_" + "idx"
	_, err := db.Exec("CREATE INDEX %s IF NOT EXISTS %s ON %s (%s) ", index, index, relation, column)
	if err != nil {
		log.Println(err)
	}
}
