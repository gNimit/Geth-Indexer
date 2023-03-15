package indexer

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "postgresql-115233-0.cloudclusters.net"
	port     = 17280
	user     = "adrem"
	password = "Mordisback@12358"
	dbname   = "multilot"
)

func connect() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Failed to establish a connection with database. %v\n", err)
	}

	return db
}
