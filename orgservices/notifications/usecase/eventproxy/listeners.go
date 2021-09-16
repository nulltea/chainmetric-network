package eventproxy

import (
	"context"
	"fmt"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
)

var (
	ctx       context.Context
	cancelAll context.CancelFunc
	cancelMap map[string]context.CancelFunc
)

func init() {
	ctx, cancelAll = context.WithCancel(context.Background())
	cancelMap = make(map[string]context.CancelFunc)
}

func Include(tickets ...model.SubscriptionTicket) error {
	if err := repository.NewSubscriptionsMongo(core.MongoDB).Insert(tickets...); err != nil {
		return fmt.Errorf("failed to persist new subscribtions tickets: %w", err)
	}

	spawnListeners(tickets...)

	return nil
}

func Revoke(ids ...string) {
	for i := range ids {
		if cancel, ok := cancelMap[ids[i]]; ok {
			cancel()
		}
	}
}

func spawnListeners(tickets ...model.SubscriptionTicket) {
	for i := range tickets {
		var (
			ticket = tickets[i]
			ctxSub context.Context
			cancel context.CancelFunc
		)

		if ticket.ExpireAt != nil {
			ctxSub, cancel = context.WithDeadline(ctx, *ticket.ExpireAt)
		} else {
			ctxSub, cancel = context.WithCancel(ctx)
		}

		cancelMap[ticket.ID] = cancel

		go eventLoop(ctxSub, &ticket)
	}
}
