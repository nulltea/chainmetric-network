package identity

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/m1/go-generate-password/generator"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	"github.com/timoth-y/chainmetric-contracts/src/identity/usecase/privileges"
)

// Enroll generates msp.SigningIdentity for user and confirms one.
func Enroll(userID string, options ...EnrollmentOption) (string, error) {
	var (
		usersRepo = repository.NewUserMongo(core.MongoDB)
		vaultRepo = repository.NewIdentitiesVault(core.Vault)
		argsStub = &enrollArgs{
			UserID: userID,
		}
	)

	for i := range options {
		options[i].Apply(argsStub)
	}

	user, err := usersRepo.GetByID(argsStub.UserID)
	if err != nil {
		return "", errors.Wrap(err, "failed to found user registration")
	}

	if err = client.Enroll(user.IdentityName(),
		msp.WithSecret(user.EnrollmentID),
		msp.WithType(identityTypeByRole(argsStub.Role)),
	); err != nil {
		return "", errors.Wrap(err, "failed to enroll user")
	}

	cert, key, err := GetSigningCredentials(user)
	if err != nil {
		return "", errors.Wrap(err, "failed to get signing identity for new user")
	}

	if _, err = vaultRepo.WriteStaticSecret(user.IdentityName(), cert, key); err != nil {
		return "", err
	}

	initialPassword, passwordHash := generatePasscode()

	if err = usersRepo.UpdateByID(user.ID, map[string]interface{}{
		"confirmed": true,
		"role":      argsStub.Role,
		"expire_at": argsStub.ExpireAt,
		"passcode":  passwordHash,
	}); err != nil {
		return "", errors.Wrap(err, "failed to update user")
	}

	if err = vaultRepo.GrantAccessWithUserpass(user.IdentityName(), passwordHash); err != nil {
		return "", err
	}

	return initialPassword, nil
}

func generatePasscode() (string, string) {
	gen, _ := generator.NewWithDefault()
	password, _ := gen.Generate()
	hash := md5.Sum([]byte(*password))
	return *password, hex.EncodeToString(hash[:])
}

func identityTypeByRole(role string) string {
	if privileges.RoleHas(role, "identity.enroll") {
		return "admin"
	}

	return "client"
}
