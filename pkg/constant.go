package pkg

import (
	"github.com/spf13/viper"
)

var (
	WsUrl           string
	ProxyUrl        string
	ProxyPublicUrl  string
	ConnectionLimit int
	AliasLength     int
)

func init() {
	// Read env file
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()

	AliasLength = 16
	viper.SetDefault("PROXY_URL", "http://localhost:5982")
	viper.SetDefault("PROXY_PUBLIC_URL", viper.GetString("PROXY_URL"))
	viper.SetDefault("WS_URL", "http://localhost:15982")
	viper.SetDefault("CONNECTION_LIMIT", 1000)

	WsUrl = viper.GetString("WS_URL")
	ProxyUrl = viper.GetString("PROXY_URL")
	ProxyPublicUrl = viper.GetString("PROXY_PUBLIC_URL")
	ConnectionLimit = viper.GetInt("CONNECTION_LIMIT")
}
