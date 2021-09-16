package model

import "time"

type Subscriber struct {
	UserID        string     `bson:"user_id"`
	SubID         string     `bson:"sub_id"`
	ExpiresAt     *time.Time `bson:"expires_at"`
	ReceivedTimes int        `bson:"received_times"`
}
