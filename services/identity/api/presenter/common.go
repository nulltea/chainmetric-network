package presenter

// NewStatusResponse presents generic StatusResponse with given `status`.
func NewStatusResponse(status Status) *StatusResponse {
	return &StatusResponse{
		Status: status,
	}
}
