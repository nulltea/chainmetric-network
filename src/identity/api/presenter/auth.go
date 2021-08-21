package presenter

// NewAuthResponse presents AuthResponse for gRPC proto with given `secretToken`
// and grants access via `jwt`.
func NewAuthResponse(secretToken, jwt string) *AuthResponse {
	return &AuthResponse{
		SecretToken: secretToken,
		AccessToken: jwt,
	}
}
