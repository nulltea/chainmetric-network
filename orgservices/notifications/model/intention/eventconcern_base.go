package intention

import (
	"context"
)

type EventConcernBase struct {
	Token       string     `bson:"subscription_token"`
	Topic    EventTopic `bson:"event_topic"`
	Contract string     `bson:"source_contract"`
	Args       map[string]interface{} `bson:"args"`
}

func (rv *EventConcernBase) SubscriptionToken() string {
	return rv.Token
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
