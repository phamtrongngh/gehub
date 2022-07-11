package pkg

import (
	"github.com/spf13/viper"
)

var (
	WsUrl           string
	ProxyUrl        string
	ConnectionLimit int
	AliasLength     int
	// Public
)

func init() {
	// Read env file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	AliasLength = 16
	viper.SetDefault("PROXY_URL", "http://localhost:5982")
	viper.SetDefault("WS_URL", "http://localhost:15982")
	viper.SetDefault("CONNECTION_LIMIT", 1000)

	WsUrl = viper.GetString("WS_URL")
	ProxyUrl = viper.GetString("PROXY_URL")
	ConnectionLimit = viper.GetInt("CONNECTION_LIMIT")
}
