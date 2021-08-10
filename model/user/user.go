package users

import "time"

type User struct {
	Firstname string `json:"firstname" bson:"firstname"`
	Lastname string `json:"lastname" bson:"lastname"`
	Email string `json:"email" bson:"email"`
	Role  string `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
