package identity

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/pkg/errors"
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
)

// GetSigningCredentials retrieves singing certificate and key for given `user`.
func GetSigningCredentials(user *model.User) (cert []byte, key []byte, err error) {
	var si msp.SigningIdentity

	if si, err = client.GetSigningIdentity(user.IdentityName()); err != nil {
		return nil, nil, errors.Wrap(err, "failed to get signing identity for new user")
	}

	cert = si.PublicVersion().EnrollmentCertificate()
	key, _ = si.PrivateKey().Bytes()

	return
}
