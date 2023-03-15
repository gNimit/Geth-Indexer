package cli

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DbOptions struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	DBNAME   string
}

type ApiOptions struct {
	ETHERSCAN_URL string
	INFURA_URL    string
}

type Options struct {
	DB    DbOptions
	API   ApiOptions
	FLAGS FlagOptions
}

func Run() *Options {

	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get path to current working directory, %v\n", err)
	}

	err = godotenv.Load(wd + "/.env")
	if err != nil {
		log.Printf("Failed to load .env file %v\n", err)
	}

	DB := DbOptions{
		HOST:     os.Getenv("host"),
		PORT:     os.Getenv("port"),
		USER:     os.Getenv("user"),
		PASSWORD: os.Getenv("password"),
		DBNAME:   os.Getenv("dbname"),
	}

	API := ApiOptions{
		ETHERSCAN_URL: os.Getenv("ETHERSCAN_URL"),
		INFURA_URL:    os.Getenv("INFURA_WS_URL"),
	}

	FLAGS := ParseArgs()

	return &Options{
		DB:    DB,
		API:   API,
		FLAGS: FLAGS,
	}
}
