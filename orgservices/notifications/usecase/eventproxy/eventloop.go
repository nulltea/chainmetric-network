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
	ctx context.Context
	cancel context.CancelFunc
	contract *fabgateway.Contract
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	contract = core.Fabric.GetContract("readings")
}

func Start() error {
	subs, err := repository.NewSubscriptionsMongo(core.MongoDB).GetAll()

	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("failed to retrieve subscription from db: %w", err)
	}

	for i := range subs {
		go eventLoop(ctx, subs[i])
	}

	return nil
}

func Stop() {
	cancel()
}

func eventLoop(ctx context.Context, sub *model.Subscription) {
	var (
		event = fmt.Sprintf("asset.%s.requirements.%s.violation", sub.AssetID, sub.Metric)
	)

	reg, notifier, err := contract.RegisterEvent(event)
	if err != nil {
		core.Logger.Errorf("failed to register event '%s': %v", event, err)

		time.Sleep(time.Minute) // TODO: algorithmic backoff
		eventLoop(ctx, sub)

		return
	}

	defer contract.Unregister(reg)

	for {
		select {
		case e := <-notifier:
			fmt.Println(e)
		case <- ctx.Done():
			return
		}
	}
}
