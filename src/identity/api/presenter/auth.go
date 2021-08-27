package presenter

import model "github.com/timoth-y/chainmetric-contracts/shared/model/user"

// NewAuthResponse presents AuthResponse for gRPC proto with given `secretToken`
// and grants access via `jwt`.
func NewAuthResponse(user *model.User, secretToken, secretPath, jwt string) *AuthResponse {
	return &AuthResponse{
		User: NewUserProto(user),
		Secret: &VaultSecret{
			Token: secretToken,
			Path: secretPath,
		},
		AccessToken: jwt,
	}
}
