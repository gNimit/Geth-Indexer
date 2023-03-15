package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

var (
	queryLotJoined   = "INSERT INTO lot_joined (lotId, token, address, size) VALUES (%d, %s, %s, %d)"
	queryLotCreated  = "INSERT INTO lot_created (lotId, tokenA, tokenBChoices, startEpoch, duration, creator, isPrivate, isChallenge) VALUES (%d, %s, %v, %d, %d, %s, %t, %t)"
	queryLotResolved = "INSERT INTO lot_resolved (lotId, size, winningToken, startPriceTokenA, startPriceTokenB, resolvePriceTokenA, resolvePriceTokenB) VALUES (%d, %d, %s, %d, %d, %d, %d)"
)

func upload(db *sql.DB, ed eventData) {
	switch ed.event {
	case "LotJoined":
		uploadJoined(db, ed.log)
	case "LotCreated":
		uploadCreated(db, ed.log)
	case "LotResolved":
		uploadResolved(db, ed.log)
	default:
		fmt.Println("Given Event not found.")
	}
}

func uploadJoined(db *sql.DB, vLog map[string]interface{}) {
	_, err := db.ExecContext(context.Background(), queryLotJoined, vLog["lotId"], vLog["token"], vLog["user"], vLog["size"])
	if err != nil {
		log.Println(err)
	}
}

func uploadCreated(db *sql.DB, vLog map[string]interface{}) {
	_, err := db.ExecContext(context.Background(), queryLotCreated, vLog["lotId"], vLog["tokenA"], vLog["tokenBChoices"], vLog["startEpoch"], vLog["duration"], vLog["creator"], vLog["isPrivate"], vLog["isChallenge"])
	if err != nil {
		log.Println(err)
	}
}

func uploadResolved(db *sql.DB, vLog map[string]interface{}) {
	_, err := db.ExecContext(context.Background(), queryLotResolved, vLog["lotId"], vLog["size"], vLog["winningToken"], vLog["startPriceTokenA"], vLog["startPriceTokenB"], vLog["resolvePriceTokenA"], vLog["resolvePriceTokenB"])
	if err != nil {
		log.Println(err)
	}
}
