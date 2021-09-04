package access

import (
	"github.com/timoth-y/chainmetric-network/orgservices/identity/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/usecase/identity"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"github.com/timoth-y/chainmetric-network/orgservices/identity/model"
)

// RequestFabricCredentialsSecret writes singing certificate and key to Vault for one-time use.
func RequestFabricCredentialsSecret(user *model.User) (string, string, error) {
	cert, key, err := identity.GetSigningCredentials(user)
	if err != nil {
		return "", "", err
	}

	path, err := repository.NewIdentitiesVault(core.Vault).WriteDynamicSecret(user.IdentityName(), cert, key)
	if err != nil {
		return "", "", err
	}

	token, err := repository.NewIdentitiesVault(core.Vault).LoginWithUserpass(user.IdentityName(), user.Passcode)
	if err != nil {
		return "", "", err
	}

	return path, token, nil
}
