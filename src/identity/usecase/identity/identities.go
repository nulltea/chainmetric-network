package identity

import (
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/infrastructure/repository"
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetSigningCredentials retrieves singing certificate and key for given `user`.
func GetSigningCredentials(user *model.User) (cert, key []byte, err error) {
	var si msp.SigningIdentity

	if si, err = client.GetSigningIdentity(user.IdentityName()); err != nil {
		return nil, nil, errors.Wrap(err, "failed to get signing identity for new user")
	}

	cert = si.PublicVersion().EnrollmentCertificate()
	key, _ = si.PrivateKey().Bytes()

	return
}

// ParseSigningCredentials parses certificate\key pair to try to identify its identity owner user.
// If certificate is valid but user doesn't exist in database -
// dummy user will be created and saved for future API interactions.
func ParseSigningCredentials(cert, key []byte) (user *model.User, err error) {
	var (
		si msp.SigningIdentity
		usersRepo = repository.NewUserMongo(core.MongoDB)
	)

	if si, err = client.CreateSigningIdentity(msp.WithCert(cert), msp.WithPrivateKey(key)); err != nil {
		return nil, errors.Wrap(err, "failed to create signing identity from cert and key pair")
	}

	id := generateUserIDFromIdentityName(si.Identifier().ID)
	if user, err = usersRepo.GetByQuery(map[string]interface{}{
		"user_id": id,
	}); err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}

		user = &model.User{
			ID: id,
			Confirmed: true,
			Role: "Admin", // TODO: need to determine role by OU
			Firstname: strings.Split(si.Identifier().ID, "@")[0],
		}

		if err = usersRepo.Upsert(*user); err != nil {
			return nil, err
		}
	}

	return user, err
}
