package identity

import (
	"fmt"

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

	si, err := client.GetSigningIdentity(user.IdentityName())
	if err != nil {
		return errors.Wrap(err, "failed to get signing identity for new user")
	}

	user.Confirmed = true
	user.Role = args.Role
	user.ExpiresAt = args.ExpireAt

	pk, _ := si.PrivateKey().Bytes()
	cert := si.PublicVersion().EnrollmentCertificate()

	fmt.Println(string(cert))
	fmt.Println(string(pk))

	if err = repo.Upsert(*user); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}
