package model

import (
	"context"
	"fmt"

	"github.com/timoth-y/chainmetric-core/models"
)

type (
	EventConcern interface {
		SubscriptionID() string
		Filter() string
		SourceContract() string
		OfTopic() EventTopic
		Context(context.Context) (context.Context, context.CancelFunc)
		Notification([]byte) (Notification, error)
	}

	EventConcernBase struct {
		ID       string     `bson:"subscription_id"`
		Topic    EventTopic `bson:"event_topic"`
		Contract string     `bson:"source_contract"`
		Args       map[string]interface{} `bson:"args"`
	}
)

func (rv *EventConcernBase) SubscriptionID() string {
	return rv.ID
}

func (rv *EventConcernBase) SourceContract() string {
	return rv.Contract
}

func (rv *EventConcernBase) Context(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

func (rv *EventConcernBase) OfTopic() EventTopic {
	return rv.Topic
}

type EventTopic string

const (
	RequirementsViolationTopic EventTopic = "requirements_violation"
)

type RequirementsViolationConcern struct {
	EventConcernBase
	Args struct{
		AssetID  string        `bson:"asset_id"`
		Metric   models.Metric `bson:"metric"`
	}
}

func (rv *RequirementsViolationConcern) Filter() string {
	return fmt.Sprintf("asset.%s.requirements.%s.violation", rv.Args.AssetID, rv.Args.Metric)
}

func (rv *RequirementsViolationConcern) Notification(data []byte) (*Notification, error) {
	return nil, nil
}
