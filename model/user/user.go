package user

import (
	"fmt"
	"time"

	"github.com/timoth-y/chainmetric-core/utils"
)

// User defines structure of user model.
//
// swagger:model
type User struct {
	// User's unique identifier
	//
	// example: f4bc94f1-3af4-4ae0-9330-19d86ca42b30
	ID string `json:"id" bson:"id"`
	// User's enrollment identifier
	//
	// example: e9f437b4-2622-47c3-8d91-921e710c7354
	EnrollmentID string `json:"enrollment_id" bson:"enrollment_id"`
	// First name of the user
	//
	// example: John
	Firstname string `json:"firstname" bson:"firstname"`
	// Last name of the user
	//
	// example: Smith
	Lastname string `json:"lastname" bson:"lastname"`
	// Email address of the user
	//
	// example: john.smith@example.com
	Email string `json:"email" bson:"email"`
	// Role of the user
	//
	// example: admin
	Role string `json:"role" bson:"role"`
	// User initial registration date and time
	//
	// example: 2021-01-02T15:04:05Z07:00
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	// User is confirmed by admin.
	Confirmed bool `json:"confirmed" bson:"confirmed"`
	// Date of user's contract expiration if defined
	//
	// example: 2021-05-02T15:04:05Z07:00
	ExpireAt *time.Time `json:"expire_at" validate:"datetime"`
}

// IdentityName forms the unique name of user's identity.
func (u *User) IdentityName() string {
	return fmt.Sprintf("user.%s", utils.Hash(u.Firstname + u.Lastname + u.ID))
}
