package auth

import (
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/identity"
)

// RequestVaultSecret writes singing certificate and key to Vault for one-time use.
func RequestVaultSecret(user *model.User) (string, error) {
	cert, key, err := identity.GetSigningCredentials(user)
	if err != nil {
		return "", err
	}

	return repository.NewIdentitiesVault(core.Vault).
		WriteDynamicSecret(user.IdentityName(), cert, key)
}
