package eventproxy

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
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

func eventLoop(ctx context.Context, concern intention.EventConcern) {
	var (
		contract           = core.Fabric.GetContract(concern.SourceContract())
		reg, notifier, err = contract.RegisterEvent(concern.Filter())
	)

	if err != nil {
		core.Logrus.WithError(err).
			WithField("filter", concern.Filter()).
			WithField("contract", concern.SourceContract()).
			Errorln("failed registering event listener")

		if err = backoff.Retry(func() error {
			reg, notifier, err = contract.RegisterEvent(concern.Filter())
			return err
		}, backoff.NewExponentialBackOff()); err != nil {
			core.Logrus.WithError(err).Errorln("backoff retry failed")
			return
		}
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
