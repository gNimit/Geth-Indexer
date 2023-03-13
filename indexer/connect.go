package eidx

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

func Connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Failed to establish a connection with database. %v\n", err)
	}
	rows, err := db.Query("CREATE TABLE LOT_JOINED ();")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rows)
	return db
}
