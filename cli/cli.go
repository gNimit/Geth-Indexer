// Package cli is used for parsing user input and reading config data.
package cli

import (
	"github.com/spf13/viper"
	"log"
)

func Run() *Config {
	//setting the name of config file
	viper.SetConfigFile("config.yml")
	//setting the path of config file
	viper.AddConfigPath(".")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file, %v\n", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Printf("Failed to unmarshall config file into the struct, %v\n", err)
	}
	config.Flag = ParseFlags()
	return &config
}
