package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/timoth-y/chainmetric-contracts/model/response"
	"github.com/timoth-y/chainmetric-contracts/shared"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"
)

// EventSocketSubscriptionTicket defines subscription ticket of event socket for metric readings.
type EventSocketSubscriptionTicket struct {
	assetID string
	metric models.Metric
	expiry time.Time
}

// BindToEventSocket creates EventSocketSubscriptionTicket for connected party,
// so that each posted models.MetricReadings record, which satisfies client request for given `assetID` and `metric`,
// would be send directly to it via event streaming subscription.
func (rc *ReadingsContract) BindToEventSocket(ctx contractapi.TransactionContextInterface, assetID, metric string) string {
	var (
		timestamp = time.Now()
		clientID, _ = ctx.GetClientIdentity().GetID()
		clientHash = utils.Hash(clientID)
		eventToken = fmt.Sprintf("%s.%s.%s", assetID, metric, clientHash)
		request = EventSocketSubscriptionTicket{
			assetID: assetID,
			metric: models.Metric(metric),
			expiry: timestamp.Add(time.Hour * 1),
		}
	)

	rc.socketTickets[eventToken] = request

	shared.Logger.Debug(fmt.Sprintf(
		"event emitter '%s' added, currently registered: %d",
		eventToken, len(rc.socketTickets),
	))

	return eventToken
}

// CloseEventSocket revokes EventSocketSubscriptionTicket for connected party
// aka close event streaming subscription for newly posted models.MetricReadings records.
func (rc *ReadingsContract) CloseEventSocket(_ contractapi.TransactionContextInterface, eventToken string) {
	delete(rc.socketTickets, eventToken)

	shared.Logger.Debug(fmt.Sprintf(
		"event emitter '%s' canceled, currently registered: %d",
		eventToken, len(rc.socketTickets),
	))
}

// sendToSocketListeners goes through EventSocketSubscriptionTicket's
// and sends response.MetricReadingsPoint for requiring ones.
func (rc *ReadingsContract) sendToSocketListeners(
	ctx contractapi.TransactionContextInterface,
	readings *models.MetricReadings,
) {
	var (
		now = time.Now()
	)

	for token, request := range rc.socketTickets {
		if now.After(request.expiry) {
			delete(rc.socketTickets, token)
			shared.Logger.Debug(fmt.Sprintf("event emitter '%s' expired, currently registered: %d",
				token, len(rc.socketTickets),
			))

			continue
		}

		if request.assetID == readings.AssetID {
			if value, ok := readings.Values[request.metric];  ok {
				point := response.MetricReadingsPoint {
					DeviceID: readings.DeviceID,
					Location: readings.Location,
					Timestamp: readings.Timestamp,
					Value: value,
				}

				if pp, err := json.Marshal(point); err == nil {
					if err = ctx.GetStub().SetEvent(token, pp); err != nil {
						shared.Logger.Error("failed to send metric readings point through event socket")
					}
				}
			}
		}
	}
}
