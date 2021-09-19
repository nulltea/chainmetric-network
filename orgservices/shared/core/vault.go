package core

import (
	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

// Vault is an instance of the Vault client for managing secrets.
var Vault *vault.Client

func initVault() {
	var (
		addr = viper.GetString("vault_address")
		err error
	)

	if Vault, err = vault.NewClient(&vault.Config{
		Address: addr,
	}); err != nil {
		Logrus.WithField("err", err).
			WithField("addr", addr).
			Fatal("failed to initialize Vault client")
	}

	Vault.SetToken(viper.GetString("vault_token"))
}
