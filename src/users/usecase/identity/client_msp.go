package identity

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	sdk    *fabsdk.FabricSDK
	client *msp.Client
)

// Init performs initialization of the identity package.
func Init() error {
	var err error

	if sdk, err = fabsdk.New(config.FromFile(viper.GetString("api.connection_config_path"))); err != nil {
		return errors.Wrap(err, "failed to connect to the blockchain network")
	}

	if client, err = msp.New(sdk.Context(
		fabsdk.WithUser("Admin"),
		fabsdk.WithOrg(viper.GetString("api.organization")),
	)); err != nil {
		return err
	}

	return nil
}
