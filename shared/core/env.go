package core

import (
	"strings"

	"github.com/spf13/viper"
)

// initEnv configures viper from environment variables.
func initEnv() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("logging", "info")

	viper.SetDefault("chaincode.ccid", "")
	viper.SetDefault("chaincode.address", "")
	viper.SetDefault("chaincode.name", "")
	viper.SetDefault("chaincode.leveldb_enabled", true)
	viper.SetDefault("chaincode.persistence_path", "app/storage")

	viper.SetDefault("api.domain", "chainmetric.network")
	viper.SetDefault("api.wallet_path", "data/wallet")
	viper.SetDefault("api.crypto_path", "data/crypto")
	viper.SetDefault("api.connection_config_path", "connection.yaml")

	viper.SetDefault("mongo_enabled", true)
	viper.SetDefault("mongo_address", "mongodb://localhost:27017")
	viper.SetDefault("mongo_connection_timeout", "10s")
	viper.SetDefault("mongo_query_timeout", "30s")
	viper.SetDefault("mongo_auth", false)
	viper.SetDefault("mongo_username", "")
	viper.SetDefault("mongo_password", "")
	viper.SetDefault("mongo_tls", false)
	viper.SetDefault("mongo_ca_cert_path", "/data/certs/mongodb-ca-cert.pem")
	viper.SetDefault("mongo_database", "chainmetric_identity")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	_ = viper.ReadInConfig()
}
