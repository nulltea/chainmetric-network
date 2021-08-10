package user

import (
	"fmt"
	"time"
)

// User defines structure of user model.
//
// swagger:model
type User struct {
	ID           string `json:"id" bson:"id"`
	EnrollmentID string `json:"enrollment_id" bson:"enrollment_id"`
	Firstname    string `json:"firstname" bson:"firstname"`
	Lastname string `json:"lastname" bson:"lastname"`
	Email string `json:"email" bson:"email"`
	Role  string `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Confirmed bool `json:"confirmed" bson:"confirmed"`
}

// IdentityName forms the unique name of user's identity.
func (u *User) IdentityName() string {
	return fmt.Sprintf("user.%s_%s", u.Firstname, u.Lastname)
}
