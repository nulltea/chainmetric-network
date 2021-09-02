package identity

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
)

// Register performs users initial registration.
func Register(options ...RegistrationOption) (*model.User, error) {
	var (
		user = &model.User{
			Confirmed: false,
			Status: model.PendingApproval,
		}
		usersRepo = repository.NewUserMongo(core.MongoDB)
		err error
	)

	for i := range options {
		options[i].Apply(user)
	}

	user.ID = generateUserIDFromIdentityName(user.IdentityName())

	if user.EnrollmentID, err = client.Register(&msp.RegistrationRequest{
		Name: user.IdentityName(),
		Type: "client",
	}); err != nil {
		return nil, errors.Wrap(err, "failed to register user")
	}

	if err = usersRepo.Upsert(*user); err != nil {
		return nil, errors.Wrap(err, "failed to store user")
	}

	return user, nil
}

func generateUserIDFromIdentityName(username string) string {
	usernameHash := md5.Sum([]byte(username))
	return hex.EncodeToString(usernameHash[:])
}
