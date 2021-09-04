package model

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
	Confirmed bool       `json:"confirmed" bson:"confirmed"`
	Status    Status     `json:"status" bson:"status"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
}

// Status defines enumeration of possibles statuses of User.
type Status int

const (
	// PendingApproval defines Status of User that is pending approval from Admin.
	PendingApproval = iota
	// Approved defines Status of User that has already been approved by Admin but hasn't yet started using app.
	Approved
	// Declined defines Status of User that has been declined by Admin.
	Declined
	// Active defines Status of User that is Approved and started using app.
	Active
	// Expired defines Status of temporary User, account of which has been expired.
	Expired
	// Canceled defines Status of User, account of which has been canceled by Admin.
	Canceled
)

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

// Encode serializer user to JSON.
func (u *User) Encode() string {
	return utils.MustEncode(u)
}

// Decode deserialize user from given JSON `payload`.
func (u User) Decode(payload string) *User {
	if err := json.Unmarshal([]byte(payload), &u); err != nil {
		panic(err)
	}
	return &u
}
