package audience

import "time"

type SubscriptionTicket struct {
	UserID        string     `bson:"user_id"`
	ConcernToken  string     `bson:"concern_token"`
	ExpiresAt     *time.Time `bson:"expires_at"`
	ReceivedTimes int        `bson:"received_times"`
}
