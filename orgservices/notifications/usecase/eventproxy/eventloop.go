package eventproxy

import (
	"context"
	"fmt"
	"time"

	fabgateway "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/infrastructure/repository"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	contract  *fabgateway.Contract
)

func init() {
	contract = core.Fabric.GetContract("readings")
}

func Start() error {
	tickets, err := repository.NewSubscriptionsMongo(core.MongoDB).GetAll()

	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("failed to retrieve subscription from db: %w", err)
	}

	spawnListeners(tickets...)
	spawnReceivers()

	return nil
}

func Stop() {
	cancelAll()
}

func eventLoop(ctx context.Context, ticket *model.SubscriptionTicket) {
	var filter = fmt.Sprintf("asset.%s.requirements.%s.violation", ticket.AssetID, ticket.Metric)

	reg, notifier, err := contract.RegisterEvent(filter)
	if err != nil {
		core.Logger.Errorf("failed to register filter '%s': %v", filter, err)

		time.Sleep(time.Minute) // TODO: algorithmic backoff
		eventLoop(ctx, ticket)

		return
	}

	defer contract.Unregister(reg)

	for {
		select {
		case e := <-notifier:
			eventsPipe <- event{
				ticket: ticket,
				payload: e.Payload,
			}
		case <- ctx.Done():
			return
		}
	}
}
