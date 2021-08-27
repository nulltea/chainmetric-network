package identity

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
)

// Enroll generates msp.SigningIdentity for user and confirms one.
func Enroll(userID string, options ...EnrollmentOption) error {
	var (
		repo = repository.NewUserMongo(core.MongoDB)
		args = &enrollArgs{
			UserID: userID,
		}
	)

	for i := range options {
		options[i].Apply(args)
	}

	user, err := repo.GetByID(args.UserID)
	if err != nil {
		return errors.Wrap(err, "failed to found user registration")
	}

	if err = client.Enroll(user.IdentityName(), msp.WithSecret(user.EnrollmentSecret)); err != nil {
		return errors.Wrap(err, "failed to enroll user")
	}

	cert, key, err := GetSigningCredentials(user)
	if err != nil {
		return errors.Wrap(err, "failed to get signing identity for new user")
	}

	if _, err = repository.NewIdentitiesVault(core.Vault).
		WriteStaticSecret(user.IdentityName(), cert, key); err != nil {
		return err
	}

	if err = repo.UpdateByID(user.ID, map[string]interface{}{
		"confirmed": true,
		"role": args.Role,
		"expire_at": args.ExpireAt,
	}); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}
