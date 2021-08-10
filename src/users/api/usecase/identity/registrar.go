package identity

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
	model "github.com/timoth-y/chainmetric-contracts/model/user"
	"github.com/timoth-y/chainmetric-contracts/shared"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/infrastructure/repository"
)

// Register performs users initial registration.
func Register(request model.RegistrationRequest) (*model.User, error) {
	var (
		user = &model.User{
			Firstname: request.Firstname,
			Lastname:  request.Lastname,
			Email:     request.Email,
		}
		err error
	)

	if user.EnrollmentID, err = client.Register(&msp.RegistrationRequest{
		Name: user.IdentityName(),
		Type: "user",
	}); err != nil {
		return nil, errors.Wrap(err, "failed to register user")
	}

	if err := repository.NewUserMongo(shared.MongoDB).Store(*user); err != nil {
		return nil, errors.Wrap(err, "failed to store user")
	}

	return user, nil
}
