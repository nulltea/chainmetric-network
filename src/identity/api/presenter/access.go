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

// NewCertificateAuthResponse presents CertificateAuthResponse for gRPC proto with given `user`
// and grants access via `jwt`.
func NewCertificateAuthResponse(user *model.User, jwt string) *CertificateAuthResponse {
	return &CertificateAuthResponse{
		User: NewUserProto(user),
		ApiAccessToken: jwt,
	}
}
