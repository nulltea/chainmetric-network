package audience

import "time"

type SubscriptionTicket struct {
	Topic         string     `bson:"topic"`
	ConcernHash   string     `bson:"concern_hash"`
	UserToken     string     `bson:"user_token"`
	ExpiresAt     *time.Time `bson:"expires_at"`
	ReceivedTimes int        `bson:"received_times"`
}
