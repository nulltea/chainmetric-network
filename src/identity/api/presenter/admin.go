package presenter

// NewEnrollUserResponse presents EnrollmentResponse for gRPC proto for given `initialPassword`.
func NewEnrollUserResponse(initialPassword string) *EnrollUserResponse {
	return &EnrollUserResponse{
		InitialPassword: initialPassword,
	}
}
