package indexer

import (
	"database/sql"
	"fmt"
	"log"

	"eventIndexer.com/cli"
	_ "github.com/lib/pq"
)

func connect(dbOps cli.DbOptions) *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", dbOps.HOST, dbOps.PORT, dbOps.USER, dbOps.PASSWORD, dbOps.DBNAME)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Failed to establish a connection with database. %v\n", err)
	}

	return db
}
