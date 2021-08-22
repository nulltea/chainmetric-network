package presenter

import model "github.com/timoth-y/chainmetric-contracts/shared/model/user"

// NewRegistrationResponse presents RegistrationResponse for gRPC proto for given `user`,
// and grants access via `jwt`.
func NewRegistrationResponse(user *model.User, jwt string) *RegistrationResponse {
	return &RegistrationResponse{
		User: NewUserProto(user),
		AccessToken: jwt,
	}
}