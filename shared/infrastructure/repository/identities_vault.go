package repository

import (
	"encoding/base64"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
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
		path = fmt.Sprintf("identity/%s/crypto", id)
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
func (r *IdentitiesVault) WriteDynamicSecret(id string, certificate, key []byte) (string, string, error) {
	var (
		path = fmt.Sprintf("auth/%s/login", id)
		data = map[string]interface{}{
			"certificate": base64.StdEncoding.EncodeToString(certificate),
			"signing_key": base64.StdEncoding.EncodeToString(key),
		}
	)

	secret, err := r.client.Logical().Write(path, data)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to write identity secret to Vault")
	}

	token, err := secret.TokenID()

	return path, token, err
}
