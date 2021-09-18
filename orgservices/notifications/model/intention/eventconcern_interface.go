package intention

import (
	"context"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/audience"
)

type (
	EventConcern interface {
		Hash() string
		Topic() string
		Filter() string
		SourceContract() string
		OfKind() EventKind
		Context(context.Context) (context.Context, context.CancelFunc)
		NotificationWith([]byte) (*audience.Notification, error)
		IsEqual(EventConcern) bool
	}

	EventKind string
)
