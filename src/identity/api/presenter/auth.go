package presenter

// NewAuthResponse presents AuthResponse for gRPC proto with given `token`.
func NewAuthResponse(token string) *AuthResponse {
	return &AuthResponse{
		SecretToken: token,
	}
}
