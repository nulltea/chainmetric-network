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
	viper.SetDefault("channel", "supply-channel")
	viper.SetDefault("logging", "info")

	viper.SetDefault("wallet_path", "/app/data")
	viper.SetDefault("crypto_path", "/app/crypto")
	viper.SetDefault("connection_config_path", "/app/config/fabric/connection.yaml")
	viper.SetDefault("fabric_cert", "/app/certs/fabric/admin.crt")
	viper.SetDefault("fabric_key", "/app/certs/fabric/admin.key")
	viper.SetDefault("grpc_tls_cert", "/app/certs/grpc/tls.crt")
	viper.SetDefault("grpc_tls_key", "/app/certs/grpc/tls.key")
	viper.SetDefault("jwt_expiration", 0)
	viper.SetDefault("jwt_signing_key", "/app/certs/jwt/jwt-key.pem")
	viper.SetDefault("jwt_signing_cert", "/app/certs/jwt/jwt-cert.pem")
	viper.SetDefault("firebase_enabled", false)
	viper.SetDefault("firebase_credentials", "/app/certs/firebase/firebase_credentials.json")
	viper.SetDefault("privileges_config","/app/config/privileges")

	viper.SetDefault("notifications.events_buffer_size", 1000)
	viper.SetDefault("notifications.event_receivers_count", 10)

	viper.SetDefault("mongo_enabled", true)
	viper.SetDefault("mongo_address", "mongodb://localhost:27018")
	viper.SetDefault("mongo_connection_timeout", "10s")
	viper.SetDefault("mongo_query_timeout", "30s")
	viper.SetDefault("mongo_auth", false)
	viper.SetDefault("mongo_username", "chainmetric_admin")
	viper.SetDefault("mongo_password", "")
	viper.SetDefault("mongo_tls", false)
	viper.SetDefault("mongo_ca_cert_path", "/app/certs/mongodb/ca-cert.pem")
	viper.SetDefault("mongo_database", "chainmetric_db")

	viper.SetDefault("vault_address", "https://vault.infra.chainmetric.network:443")
	viper.SetDefault("vault_token", "")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	_ = viper.ReadInConfig()
}
