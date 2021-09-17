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

func Include(concerns ...model.EventConcern) error {
	if err := repository.NewEventConcernsMongo(core.MongoDB).Insert(concerns...); err != nil {
		return fmt.Errorf("failed to persist new event concers: %w", err)
	}

	spawnListeners(concerns...)

	return nil
}

func Revoke(ids ...string) error {
	if err := repository.NewEventConcernsMongo(core.MongoDB).DeleteByIDs(ids...); err != nil {
		return fmt.Errorf("failed to delete event concers: %w", err)
	}

	for i := range ids {
		if cancel, ok := cancelMap[ids[i]]; ok {
			cancel()
		}
	}

	return nil
}

func spawnListeners(concerns ...model.EventConcern) {
	for _, current := range concerns {
		ctx, cancel := current.Context(ctx)
		cancelMap[current.SubscriptionID()] = cancel

		go eventLoop(ctx, current)
	}
}
