package user

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/utils"
)

// User defines structure of user model.
type User struct {
	ID           string     `json:"id" bson:"id"`
	EnrollmentID string     `json:"enrollment_id" bson:"enrollment_id"`
	Firstname    string     `json:"firstname" bson:"firstname"`
	Lastname     string     `json:"lastname" bson:"lastname"`
	Email        string     `json:"email" bson:"email"`
	Role         string     `json:"role,omitempty" bson:"role,omitempty"`
	Passcode     string     `json:"passcode" bson:"passcode"`
	CreatedAt    time.Time  `json:"created_at" bson:"created_at"`
	Confirmed    bool       `json:"confirmed" bson:"confirmed"`
	Trained      bool       `json:"trained" bson:"trained"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

// IdentityName forms the unique name of user's identity.
func (u *User) IdentityName() string {
	usernameParts := []string {
		strings.ToLower(u.Firstname),
	}

	if len(u.Lastname) != 0 {
		usernameParts = append(usernameParts, strings.ToLower(u.Lastname))
	}

	return fmt.Sprintf("%s@%s.org.%s",
		strings.Join(usernameParts, "."),
		viper.GetString("organization"),
		viper.GetString("domain"),
	)
}

func (u *User) Encode() string {
	return utils.MustEncode(u)
}

func (u User) Decode(payload string) *User {
	if err := json.Unmarshal([]byte(payload), &u); err != nil {
		panic(err)
	}
	return &u
}
