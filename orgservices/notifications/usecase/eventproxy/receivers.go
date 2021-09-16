package eventproxy

import (
	"context"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model"
)

type event struct {
	ticket *model.SubscriptionTicket
	payload []byte
}

var (
	eventsPipe chan event
)

func init() {
	eventsPipe = make(chan event, viper.GetInt("api.notifications.events_buffer_size"))
}

func spawnReceivers() {
	var receiverCount = viper.GetInt("api.notifications.event_receivers_count")

	for i := 0; i < receiverCount; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case e := <-eventsPipe:
					_ = e // TODO: redirect to subscribed users
				case <-ctx.Done():
					return
				}
			}
		}(ctx)
	}
}
