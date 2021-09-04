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

	viper.SetDefault("api.wallet_path", "./data")
	viper.SetDefault("api.crypto_path", "./crypto")
	viper.SetDefault("api.connection_config_path", "config/connection.yaml")
	viper.SetDefault("api.grpc_tls_cert", "./certs/grpc/tls.crt")
	viper.SetDefault("api.grpc_tls_key", "./certs/grpc/tls.key")
	viper.SetDefault("api.jwt_expiration", 0)
	viper.SetDefault("api.jwt_signing_key", "./certs/jwt/jwt-key.pem")
	viper.SetDefault("api.jwt_signing_cert", "./certs/jwt/jwt-cert.pem")

	viper.SetDefault("mongo_enabled", true)
	viper.SetDefault("mongo_address", "mongodb://localhost:27017")
	viper.SetDefault("mongo_connection_timeout", "10s")
	viper.SetDefault("mongo_query_timeout", "30s")
	viper.SetDefault("mongo_auth", false)
	viper.SetDefault("mongo_username", "chainmetric_admin")
	viper.SetDefault("mongo_password", "")
	viper.SetDefault("mongo_tls", false)
	viper.SetDefault("mongo_ca_cert_path", "certs/mongodb/ca-cert.pem")
	viper.SetDefault("mongo_database", "chainmetric_db")

	viper.SetDefault("vault_address", "https://vault.infra.chainmetric.network:443")
	viper.SetDefault("vault_token", "")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	_ = viper.ReadInConfig()
}
