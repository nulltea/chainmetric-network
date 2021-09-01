package presenter

import model "github.com/timoth-y/chainmetric-contracts/shared/model/user"

// NewFabricCredentialsResponse presents FabricCredentialsResponse for gRPC proto with given `secretToken`
// and grants access via `jwt`.
func NewFabricCredentialsResponse(user *model.User, secretToken, secretPath, jwt string) *FabricCredentialsResponse {
	return &FabricCredentialsResponse{
		User: NewUserProto(user),
		Secret: &VaultSecret{
			Token: secretToken,
			Path: secretPath,
		},
		ApiAccessToken: jwt,
	}
}
