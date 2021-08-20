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
	var (
		// domain = viper.GetString("api.domain")
		// org = viper.GetString("api.organization")
		err error
	)

	if sdk, err = fabsdk.New(config.FromFile(viper.GetString("api.connection_config_path"))); err != nil {
		return errors.Wrap(err, "failed to connect to the blockchain network")
	}

	if client, err = msp.New(
		sdk.Context(),
		msp.WithOrg(viper.GetString("api.organization")),
	); err != nil {
		return err
	}

	// adminID := fmt.Sprintf("Admin@%s.org.%s", org, domain )
	//
	// if ir, err := client.GetIdentity(adminID); err != nil || ir == nil {
	// 	if err := client.Enroll(adminID,
	// 		msp.WithSecret("adminpsw"),
	// 	); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
