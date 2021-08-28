package user

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/timoth-y/chainmetric-core/utils"
)

// User defines structure of user model.
type User struct {
	// User's unique identifier
	ID string `json:"id" bson:"id" example:"f4bc94f1-3af4-4ae0-9330-19d86ca42b30"`
	// First name of the user
	Firstname string `json:"firstname" bson:"firstname" example:"John"`
	// Last name of the user
	Lastname string `json:"lastname" bson:"lastname" example:"Smith"`
	// Email address of the user
	Email string `json:"email" bson:"email" example:"john.smith@example.com"`
	// Role of the user
	Role string `json:"role,omitempty" bson:"role,omitempty" example:"admin"`
	// User initial registration date and time
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2021-01-02T15:04:05Z07:00"`
	// User is confirmed by admin.
	Confirmed bool `json:"confirmed" bson:"confirmed"`
	// Date of user's contract expiration if defined
	ExpiresAt *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty" example:"2021-05-02T15:04:05Z07:00"`

	EnrollmentSecret string `json:"-" bson:"enrollment_id"`
	Passcode         string `json:"-" bson:"passcode"`
}

// IdentityName forms the unique name of user's identity.
func (u *User) IdentityName() string {
	return fmt.Sprintf("%s.%s", strings.ToLower(u.Firstname), strings.ToLower(u.Lastname))
}

func (u *User) Encode() string {
	return utils.MustEncode(u)
}

func (u User) Decode(payload string) *User {
	json.Unmarshal([]byte(payload), &u)
	return &u
}
