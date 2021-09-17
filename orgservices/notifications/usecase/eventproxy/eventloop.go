package eventproxy

import (
	"context"
	"fmt"
	"time"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"go.mongodb.org/mongo-driver/mongo"
)

func Start() error {
	concerns, err := repository.NewEventConcernsMongo(core.MongoDB).GetAll()

	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("failed to retrieve event concers from db: %w", err)
	}

	spawnListeners(concerns...)
	spawnReceivers()

	return nil
}

func Stop() {
	cancelAll()
}

func eventLoop(ctx context.Context, concern model.EventConcern) {
	var (
		contract           = core.Fabric.GetContract(concern.SourceContract())
		reg, notifier, err = contract.RegisterEvent(concern.Filter())
	)

	if err != nil {
		core.Logger.Errorf("failed to register filter '%s': %v", concern.Filter(), err)

		time.Sleep(time.Minute) // TODO: algorithmic backoff
		eventLoop(ctx, concern)

		return
	}

	defer contract.Unregister(reg)

	for {
		select {
		case e := <-notifier:
			eventsPipe <- event{
				EventConcern: concern,
				payload: e.Payload,
			}
		case <- ctx.Done():
			return
		}
	}
}
