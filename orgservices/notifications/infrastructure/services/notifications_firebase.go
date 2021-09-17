package services

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model"
)

// NotificationsFirebase defines service for sending notification via Firebase Cloud Messaging protocol.
type NotificationsFirebase struct {
	app *firebase.App
}

// NewNotificationsFirebase constructs new NotificationsFirebase service instance.
func NewNotificationsFirebase(app *firebase.App) *NotificationsFirebase {
	return &NotificationsFirebase{
		app: app,
	}
}

func (nf *NotificationsFirebase) Push(userID string, n model.Notification) error {
	client, err := nf.app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("failed to use message client: %w", err)
	}

	if _, err = client.Send(context.Background(), &messaging.Message{
		Topic: string(n.Topic),
		Notification: &messaging.Notification{
			Title: n.Caption,
			Body: n.Description,
		},
		Data: map[string]string{
			"data": utils.MustEncode(n.Data),
		},
	}); err != nil {
		return fmt.Errorf("failed to send %s notification: %w", n.Topic, err)
	}

	return nil
}

