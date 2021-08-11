package identity

import (
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
	model "github.com/timoth-y/chainmetric-contracts/model/user"
	"github.com/timoth-y/chainmetric-contracts/shared"
	"github.com/timoth-y/chainmetric-contracts/src/users/infrastructure/repository"
)

// Register performs users initial registration.
func Register(request model.RegistrationRequest) (*model.User, error) {
	var (
		user = &model.User{
			ID:        uuid.NewString(),
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

	if err := repository.NewUserMongo(shared.MongoDB).Upsert(*user); err != nil {
		return nil, errors.Wrap(err, "failed to store user")
	}

	return user, nil
}

// Enroll generates msp.SigningIdentity for user and confirms one.
func Enroll(req model.EnrollmentRequest) error {
	var (
		repo = repository.NewUserMongo(shared.MongoDB)
	)

	user, err := repo.GetByID(req.UserID)
	if err != nil {
		return errors.Wrap(err, "failed to found user registration")
	}

	if err = client.Enroll(user.EnrollmentID); err != nil {
		return errors.Wrap(err, "failed to enroll user")
	}

	si, err := client.GetSigningIdentity(user.IdentityName())
	if err != nil {
		return errors.Wrap(err, "failed to get signing identity for new user")
	}

	_ = si
	user.Confirmed = true
	user.Role = req.Role
	user.ExpireAt = req.ExpireAt

	if err = repo.Upsert(*user); err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	return nil
}
