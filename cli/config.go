// Package cli is used for parsing user input and reading config data.
package cli

// Config is Configuration Data Exported also contained user input from flags.
type Config struct {
	Query    FlagOptions
	Database DatabaseConfig
	API      APIConfig
}

// DatabaseConfig is database Credentials of user
type DatabaseConfig struct {
	DBHost     string `mapstructure:"host"`
	DBPort     int    `mapstructure:"port"`
	DBUser     string `mapstructure:"user"`
	DBPassword string `mapstructure:"password"`
	DBName     string `mapstructure:"name"`
}

// APIConfig is APIs URL for connecting to Etherscan to fetch ABI of contract
// and establish connection with a ethereum node local or remote.
type APIConfig struct {
	EtherScanAPI string `mapstructure:"etherscan"`
	EthNodeURL   string `mapstructure:"ethnode"`
}
