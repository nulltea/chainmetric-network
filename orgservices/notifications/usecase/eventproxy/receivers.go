package eventproxy

import (
	"context"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/infrastructure/services"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
)

type event struct {
	intention.EventConcern
	payload []byte
}

var (
	eventsPipe chan event
	fcmService *services.NotificationsFirebase
)

func init() {
	eventsPipe = make(chan event, viper.GetInt("api.notifications.events_buffer_size"))
	fcmService = services.NewNotificationsFirebase(core.Firebase)
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
			WithField("topic", e.OfTopic()).
			WithField("sub_token", e.SubscriptionToken()).
			Errorf("failed to form notification")

	} else if err = fcmService.Push(notification); err != nil {
		core.Logrus.WithError(err).
			WithField("topic", e.Filter()).
			WithField("sub_token", e.SubscriptionToken()).
			Errorf("failed to push notification via FCM")
	}
}
