package core

import (
	"io/ioutil"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/spf13/viper"
)

var (
	wallet  *gateway.Wallet
	Fabric *gateway.Network
)

func initFabric() {
	var (
		signCert = viper.GetString("fabric_cert")
		signKey = viper.GetString("fabric_key")
		channel = viper.GetString("channel")
	)

	wallet = gateway.NewInMemoryWallet()
	identity := gateway.NewX509Identity(viper.GetString("organization"), "", "")

	if payload, err := ioutil.ReadFile(signCert); err != nil {
		Logger.Fatalf("failed to read fabric identity certificate on path %s: %v", signCert, err)
	} else {
		identity.Credentials.Certificate = string(payload)
	}

	if payload, err := ioutil.ReadFile(signKey); err != nil {
		Logger.Fatalf("failed to read fabric identity key on path %s: %v",
			signKey, err)
	} else {
		identity.Credentials.Key = string(payload)
	}

	wallet.Put("service", identity)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(viper.GetString("connection_config_path"))),
		gateway.WithIdentity(wallet, "service"),
	)

	if err != nil {
		Logger.Fatalf("failed to connect to Fabric gateway: %v", err)
	}

	if Fabric, err = gw.GetNetwork(channel); err != nil {
		Logger.Fatalf("failed to init to Fabric network client on channel %s: %v", channel, err)
	}
}
