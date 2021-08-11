package user

import "time"

// RegistrationRequest defines structure of the request to register new User.
type RegistrationRequest struct {
	// First name of the user
	Firstname string `json:"firstname" validate:"required" example:"John"`
	// Last name of the user
	Lastname string `json:"lastname" validate:"required" example:"Smith"`
	// Email address of the user
	Email string `json:"email" validate:"required,email" example:"john.smith@example.com"`
}

// EnrollmentRequest defines structure of the request to enroll new User.
type EnrollmentRequest struct {
	// User's unique identifier
	UserID string `json:"user_id" validate:"required,uuid" example:"f4bc94f1-3af4-4ae0-9330-19d86ca42b30"`
	// Role of the user
	Role   string `json:"role" validate:"required" example:"admin"`
	// Date of user's contract expiration if defined
	ExpireAt *time.Time `json:"expire_at" validate:"datetime" example:"2006-01-02T15:04:05Z07:00"`
}
