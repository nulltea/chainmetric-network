package audience

import "time"

type SubscriptionTicket struct {
	UserToken        string  `bson:"user_token"`
	ConcernToken  string     `bson:"concern_token"`
	ExpiresAt     *time.Time `bson:"expires_at"`
	ReceivedTimes int        `bson:"received_times"`
}
