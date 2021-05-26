package shared

import (
	"github.com/spf13/viper"
)

// initEnv configures viper from environment variables.
func initEnv() {
	viper.SetEnvPrefix("CHAINCODE")
	viper.AutomaticEnv()

	viper.SetDefault("ccid", "")
	viper.SetDefault("address", "")
	viper.SetDefault("name", "")

	viper.SetDefault("logging", "info")

	viper.SetDefault("persistence_path", "./storage/")
}
