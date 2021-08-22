package identity

import (
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	"github.com/timoth-y/chainmetric-contracts/shared/model/user"
)

// Register performs users initial registration.
func Register(options ...RegistrationOption) (*user.User, error) {
	var (
		user = &user.User{
			ID: uuid.NewString(),
		}
		err error
	)

	for i := range options {
		options[i].Apply(user)
	}

	if user.EnrollmentSecret, err = client.Register(&msp.RegistrationRequest{
		Name: user.IdentityName(),
		Type: "client",
	}); err != nil {
		return nil, errors.Wrap(err, "failed to register user")
	}

	if err := repository.NewUserMongo(core.MongoDB).Upsert(*user); err != nil {
		return nil, errors.Wrap(err, "failed to store user")
	}

	return user, nil
}
