package model

import (
	"github.com/timoth-y/chainmetric-core/models"
)

type Subscription struct {
	ID      string        `bson:"id"`
	AssetID string        `bson:"asset_id"`
	Metric  models.Metric `bson:"metric"`

	Subscribers []Subscriber
}


