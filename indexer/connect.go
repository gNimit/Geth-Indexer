package indexer

import (
	"database/sql"
	"eventIndexer.com/cli"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(dbOpts cli.DatabaseConfig) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", dbOpts.DBHost, dbOpts.DBPort, dbOpts.DBUser, dbOpts.DBPassword, dbOpts.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
