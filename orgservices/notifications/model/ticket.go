package model

import (
	"time"

	"github.com/timoth-y/chainmetric-core/models"
)

type SubscriptionTicket struct {
	ID       string        `bson:"id"`
	AssetID  string        `bson:"asset_id"`
	Metric   models.Metric `bson:"metric"`
	ExpireAt *time.Time    `bson:"expire_at"`

	Subscribers []Subscriber
}


