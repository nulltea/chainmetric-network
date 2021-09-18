package eventproxy

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/infrastructure/services"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/audience"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
)

type event struct {
	intention.EventConcern
	payload []byte
}

var (
	eventsPipe chan event
	subsRepo   *repository.SubscriptionsMongo
	fcmService *services.NotificationsFirebase
)

func init() {
	eventsPipe = make(chan event, viper.GetInt("api.notifications.events_buffer_size"))
	fcmService = services.NewNotificationsFirebase(core.Firebase)
	subsRepo = repository.NewSubscriptionsMongo(core.MongoDB)
}

// ForwardEvents subscribes user with given `token` to events defined by given `concerns`,
// and includes them for redirecting via Include.
func ForwardEvents(userToken string, concerns ...intention.EventConcern) error {
	var (
		topics  = make([]string, len(concerns))
		tickets = make([]audience.SubscriptionTicket, len(concerns))
	)

	for i := range concerns {
		topics[i] = concerns[i].Topic()
		tickets[i] = audience.SubscriptionTicket{
			Topic: concerns[i].Topic(),
			ConcernHash: concerns[i].Hash(),
			UserToken: userToken,
		}
	}

	if err := fcmService.SubscribeToTopics(userToken, topics...); err != nil {
		return fmt.Errorf("failed subscribing user token '%s' to topics (%s): %w",
			userToken, strings.Join(topics, ","), err,
		)
	}

	if err := subsRepo.Insert(tickets...); err != nil {
		return fmt.Errorf("failed log subscribtion tickets for user token '%s': %w", userToken, err)
	}

	return Include(concerns...)
}

// CancelForwarding ...
func CancelForwarding(userToken string, topics ...string) error {
	if err := subsRepo.DeleteByTopicsForUser(userToken, topics...); err != nil {
		return fmt.Errorf("failed delete subscribtion tickets: %w", err)
	}

	if err := fcmService.UnsubscribeFromTopics(userToken, topics...); err != nil {
		return fmt.Errorf("failed unsubscribing user token '%s' from topics (%s): %w",
			userToken, strings.Join(topics, ","), err,
		)
	}

	remainMap, err := subsRepo.CountByTopics(userToken, topics...)
	if err != nil {
		for topic, remained := range remainMap {
			if remained == 0 {
				Revoke(topic) //
			}
		}
	}

	return nil
}

func spawnReceivers() {
	var receiverCount = viper.GetInt("api.notifications.event_receivers_count")

	for i := 0; i < receiverCount; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case e := <-eventsPipe:
					forwardEvent(e)
				case <-ctx.Done():
					return
				}
			}
		}(ctx)
	}
}

func forwardEvent(e event) {
	notification, err := e.NotificationWith(e.payload)

	if err != nil {
		core.Logrus.WithError(err).
			WithField("kind", e.OfKind()).
			WithField("topic", e.Topic()).
			Errorf("failed to form notification")

	} else if err = fcmService.Push(e.Topic(), notification); err != nil {
		core.Logrus.WithError(err).
			WithField("topic", e.Topic()).
			Errorf("failed to push notification via FCM")
	}
}
