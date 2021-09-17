package intention

import (
	"context"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/audience"
)

type (
	EventConcern interface {
		SubscriptionToken() string
		Filter() string
		SourceContract() string
		OfTopic() EventTopic
		Context(context.Context) (context.Context, context.CancelFunc)
		NotificationFor(string, []byte) (*audience.Notification, error)
	}

	EventTopic string
)
