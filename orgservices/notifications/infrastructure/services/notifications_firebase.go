package services

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/audience"
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

func (nf *NotificationsFirebase) Push(nft *audience.Notification) error {
	client, err := nf.app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("failed to use message client: %w", err)
	}

	if _, err = client.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: nft.Caption,
			Body: nft.Description,
		},
		Data: map[string]string{
			"data": utils.MustEncode(nft.Data),
		},
		Topic: nft.Topic,
	}); err != nil {
		return fmt.Errorf("failed to send %s notification: %w", nft.Topic, err)
	}

	return nil
}

func (nf *NotificationsFirebase) SubscribeToTopic(topic string, userTokens ...string) error {
	client, err := nf.app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("failed to use message client: %w", err)
	}

	if _, err = client.SubscribeToTopic(context.Background(), userTokens, topic); err != nil {
		return fmt.Errorf("failed subscribing users to %s topic: %w", topic, err)
	}

	return nil
}

func (nf *NotificationsFirebase) UnsubscribeFromTopic(topic string, userTokens ...string) error {
	client, err := nf.app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("failed to use message client: %w", err)
	}

	if _, err = client.UnsubscribeFromTopic(context.Background(), userTokens, topic); err != nil {
		return fmt.Errorf("failed subscribing users to %s topic: %w", topic, err)
	}

	return nil
}
