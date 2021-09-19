package core

import (
	"strings"

	"github.com/spf13/viper"
)

// initConfig configures viper from environment variables and yaml files.
func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("domain", "chainmetric.network")
	viper.SetDefault("organization", "chipa-inu")
	viper.SetDefault("logging", "info")

	viper.SetDefault("chaincode.ccid", "")
	viper.SetDefault("chaincode.address", "")
	viper.SetDefault("chaincode.name", "")
	viper.SetDefault("chaincode.leveldb_enabled", true)
	viper.SetDefault("chaincode.persistence_path", "app/storage")

	viper.SetDefault("vault_address", "https://vault.infra.chainmetric.network:443")
	viper.SetDefault("vault_token", "")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	_ = viper.ReadInConfig()
}
