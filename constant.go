package main

import "github.com/spf13/viper"

var (
	Port            string
	ConnectionLimit int
	AliasLength     int
)

func init() {
	// Read env file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	viper.SetDefault("PORT", 9999)
	viper.SetDefault("CONNECTION_LIMIT", 1000)

	AliasLength = 16
	Port = viper.GetString("PORT")
	ConnectionLimit = viper.GetInt("CONNECTION_LIMIT")
}
