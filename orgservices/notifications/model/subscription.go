package model

import "time"

type SubscriptionTicket struct {
	UserID    string     `bson:"user_id"`
	TicketID  string     `bson:"sub_id"`
	ExpiresAt *time.Time `bson:"expires_at"`
	ReceivedTimes int        `bson:"received_times"`
}
