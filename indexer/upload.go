package indexer

import (
	"database/sql"
	"log"
)

func upload(db *sql.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
