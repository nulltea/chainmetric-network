package user

import "time"

// RegistrationRequest defines structure of the request to register new User.
//
// swagger:model
type RegistrationRequest struct {
	// First name of the user
	//
	// example: John
	// required: true
	Firstname string `json:"firstname" validate:"required"`
	// Last name of the user
	//
	// example: Smith
	// required: true
	Lastname string `json:"lastname" validate:"required"`
	// Email address of the user
	//
	// example: john.smith@example.com
	// required: true
	Email string `json:"email" validate:"required,email"`
}

// EnrollmentRequest defines structure of the request to enroll new User.
//
// swagger:model
type EnrollmentRequest struct {
	// User's unique identifier
	//
	// example: f4bc94f1-3af4-4ae0-9330-19d86ca42b30
	// required: true
	UserID string `json:"user_id" validate:"required,uuid"`

	// Role of the user
	//
	// example: admin
	// required: true
	Role   string `json:"role" validate:"required"`

	// Date of user's contract expiration if defined.
	//
	// example: 2006-01-02T15:04:05Z07:00
	// required: false
	ExpireAt *time.Time `json:"expire_at" validate:"datetime"`
}
