package crypto

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	wallet *gateway.Wallet
)

func Init() error {
	var (
		path = viper.GetString("api.wallet_path")
		err  error
	)

	if wallet, err = gateway.NewFileSystemWallet(path); err != nil {
		return errors.Wrapf(err, "failed to create crypto wallet on path %s", path)
	}

	if !wallet.Exists("admin") {
		identity, err := getAdminIdentity();
		if err != nil {
			return err
		}

		if err := wallet.Put("admin", identity); err != nil {
			return errors.Wrap(err, "failed to save admin identity in wallet")
		}
	}

	return nil
}

func getAdminIdentity() (*gateway.X509Identity, error) {
	var (
		orgID    = viper.GetString("api.organization")
		orgHost = fmt.Sprintf("%s.org.%s", orgID, viper.GetString("api.domain"))
		adminUser = fmt.Sprintf("Admin@%s", orgHost)
		mspPath = filepath.Join(
			viper.GetString("api.crypto_path"),
			"peerOrganizations",
			orgHost,
			"users",
			adminUser,
			"msp",
		)
		certPath = filepath.Join(mspPath, "signcerts", fmt.Sprintf("%s-cert.pem", adminUser))
		keysPath = filepath.Join(mspPath, "keystore")
	)

	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, errors.Wrap(err, "failed to find identity certificate for Admin")
	}

	files, err := ioutil.ReadDir(keysPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find identity private key for Admin")
	} else if len(files) != 1 {
		return nil, errors.New("keystore folder should have contain one file")
	}

	keyPath := filepath.Join(keysPath, files[0].Name())
	pk, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read identity private key for Admin")
	}

	return gateway.NewX509Identity(orgID, string(cert), string(pk)), nil
}
