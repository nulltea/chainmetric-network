package repository

import (
	"encoding/base64"
	"fmt"
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// IdentitiesVault defines identity repository for Vault secret manager.
type IdentitiesVault struct {
	client *vault.Client
	usernameFormatter func(username string) string
}

// NewIdentitiesVault constructs new IdentitiesVault repository instance.
func NewIdentitiesVault(client *vault.Client) *IdentitiesVault {
	return &IdentitiesVault{
		client: client,
		usernameFormatter: func(username string) string {
			return fmt.Sprintf("%s.%s", strings.Split(username, "@")[0], viper.GetString("organization"))
		},
	}
}

// WriteStaticSecret writes signing credentials to Vault as static secret.
func (r *IdentitiesVault) WriteStaticSecret(username string, certificate, key []byte) (string, error) {
	var (
		path = fmt.Sprintf("fabric/identity/%s/%s", viper.GetString("organization"),
			r.usernameFormatter(username))
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
func (r *IdentitiesVault) WriteDynamicSecret(username string, certificate, key []byte) (string, error) {
	var (
		path = fmt.Sprintf("fabric/auth/login/%s", r.usernameFormatter(username))
		data = map[string]interface{}{
			"certificate": base64.StdEncoding.EncodeToString(certificate),
			"signing_key": base64.StdEncoding.EncodeToString(key),
		}
	)

	if _, err := r.client.Logical().Write(path, data); err != nil {
		return "", errors.Wrap(err, "failed to write identity secret to Vault")
	}

	return path, nil
}

// GrantAccessWithUserpass creates user in Vault,
// with given `username` and `password` credentials for userpass auth.
func (r *IdentitiesVault) GrantAccessWithUserpass(username, password string) error {
	var (
		usernameFormatted = r.usernameFormatter(username)
		path = fmt.Sprintf("auth/userpass/users/%s", usernameFormatted)
		data = map[string]interface{}{
			"password": password,
			"policies": strings.Join([]string{usernameFormatted, "default"}, ","),
		}
	)

	if _, err := r.client.Logical().Write(path, data); err != nil {
		return errors.Wrap(err, "failed to create access for Vault")
	}

	if err := r.client.Sys().PutPolicy(usernameFormatted, fmt.Sprintf(`path "fabric/auth/login/%s" {
	capabilities = [ "read" ]
}`, usernameFormatted)); err != nil {
		return errors.Wrapf(err, "failed to grand user access to '%s' path", path)
	}

	return nil
}

// UpdateUserpassAccess updates `password` for exciting Vault user account.
func (r *IdentitiesVault) UpdateUserpassAccess(username, password string) error {
	var (
		path = fmt.Sprintf("auth/userpass/users/%s/password", r.usernameFormatter(username))
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
		path = fmt.Sprintf("auth/userpass/login/%s", r.usernameFormatter(username))
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
