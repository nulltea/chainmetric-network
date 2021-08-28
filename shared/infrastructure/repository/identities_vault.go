package repository

import (
	"encoding/base64"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// IdentitiesVault defines identity repository for Vault secret manager.
type IdentitiesVault struct {
	client *vault.Client
}

// NewIdentitiesVault constructs new IdentitiesVault repository instance.
func NewIdentitiesVault(client *vault.Client) *IdentitiesVault {
	return &IdentitiesVault{
		client: client,
	}
}

// WriteStaticSecret writes signing credentials to Vault as static secret.
func (r *IdentitiesVault) WriteStaticSecret(id string, certificate, key []byte) (string, error) {
	var (
		path = fmt.Sprintf("fabric/identity/%s/%s", viper.GetString("organization"), id)
		data = map[string]interface{}{
			"certificate": base64.StdEncoding.EncodeToString(certificate),
			"signing_key": base64.StdEncoding.EncodeToString(key),
		}
	)

	_, err := r.client.Logical().Write(path, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to write identity secret to Vault")
	}

	return path, nil
}

// WriteDynamicSecret writes signing credentials to Vault for one-time use.
func (r *IdentitiesVault) WriteDynamicSecret(id string, certificate, key []byte) (string, error) {
	var (
		path = fmt.Sprintf("fabric/auth/login/%s", id)
		data = map[string]interface{}{
			"certificate": base64.StdEncoding.EncodeToString(certificate),
			"signing_key": base64.StdEncoding.EncodeToString(key),
		}
	)

	_, err := r.client.Logical().Write(path, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to write identity secret to Vault")
	}

	return path, err
}

// GrantAccessWithUserpass creates user in Vault,
// with given `username` and `password` credentials for userpass auth.
func (r *IdentitiesVault) GrantAccessWithUserpass(username, password string) error {
	var (
		path = fmt.Sprintf("auth/userpass/users/%s", username)
		data = map[string]interface{}{
			"password": password,
		}
	)

	if _, err := r.client.Logical().Write(path, data); err != nil {
		return errors.Wrap(err, "failed to create access for Vault")
	}

	return nil
}

// UpdateUserpassAccess updates `password` for exciting Vault user account.
func (r *IdentitiesVault) UpdateUserpassAccess(username, password string) error {
	var (
		path = fmt.Sprintf("auth/userpass/users/%s/password", username)
		data = map[string]interface{}{
			"password": password,
		}
	)

	if _, err := r.client.Logical().Write(path, data); err != nil {
		return errors.Wrap(err, "failed to update access for Vault")
	}

	return nil
}

// LoginWithUserpass requests access token for user with given `username` and `password` credentials.
func (r *IdentitiesVault) LoginWithUserpass(username, password string) (string, error) {
	var (
		path = fmt.Sprintf("auth/userpass/login/%s", username)
		data = map[string]interface{}{
			"password": password,
		}
	)

	secret, err := r.client.Logical().Write(path, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to authenticate user in Vault")
	}

	token, err := secret.TokenID()

	return token, err
}
