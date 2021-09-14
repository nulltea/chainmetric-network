package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/model/response"
)

// EventSocketSubscriptionTicket defines subscription ticket of event socket for metric readings.
type EventSocketSubscriptionTicket struct {
	AssetID string        `json:"asset_id"`
	Metric  models.Metric `json:"metric"`
	Expiry  time.Time     `json:"expiry"`
}

// BindToEventSocket creates EventSocketSubscriptionTicket for connected party,
// so that each posted models.MetricReadings record, which satisfies client request for given `assetID` and `metric`,
// would be sent directly to it via event streaming subscription.
func (rc *ReadingsContract) BindToEventSocket(ctx contractapi.TransactionContextInterface, assetID, metric string) string {
	var (
		timestamp   = time.Now()
		clientID, _ = ctx.GetClientIdentity().GetID()
		clientHash  = utils.Hash(clientID)
		eventToken  = fmt.Sprintf("%s.%s.%s", assetID, metric, clientHash)
		ticket      = EventSocketSubscriptionTicket{
			AssetID: assetID,
			Metric:  models.Metric(metric),
			Expiry:  timestamp.Add(time.Hour * 1),
		}
	)

	rc.socketTickets[eventToken] = ticket
	backupTicket(eventToken, ticket)

	core.Logger.Debug(fmt.Sprintf(
		"event emitter '%s' added, currently registered: %d",
		eventToken, len(rc.socketTickets),
	))

	return eventToken
}

// CloseEventSocket revokes EventSocketSubscriptionTicket for connected party
// aka close event streaming subscription for newly posted models.MetricReadings records.
func (rc *ReadingsContract) CloseEventSocket(_ contractapi.TransactionContextInterface, eventToken string) {
	delete(rc.socketTickets, eventToken)
	dropTicketBackup(eventToken)

	core.Logger.Debug(fmt.Sprintf(
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
	var now = time.Now()

	for token, ticket := range rc.socketTickets {
		if now.After(ticket.Expiry) {
			dropTicketBackup(token)
			core.Logger.Debug(fmt.Sprintf("event emitter '%s' expired, currently registered: %d",
				token, len(rc.socketTickets),
			))

			continue
		}

		if ticket.AssetID == readings.AssetID {
			if value, ok := readings.Values[ticket.Metric];  ok {
				point := response.MetricReadingsPoint{
					DeviceID: readings.DeviceID,
					Location: readings.Location,
					Timestamp: readings.Timestamp,
					Value: value,
				}

				if pp, err := json.Marshal(point); err == nil {
					if err = ctx.GetStub().SetEvent(token, pp); err != nil {
						core.Logger.Error("failed to send metric readings point through event socket")
					}
				}
			}
		}
	}
}

func (rc *ReadingsContract) recoverEventTicketsFromBackup() {
	var (
		now = time.Now()
		prefix = []byte(utils.FormCompositeKey("ticket"))
		iter = core.LevelDB.NewIterator(util.BytesPrefix(prefix), nil)
	)

	for iter.Next() {
		var (
			ticket EventSocketSubscriptionTicket
			key    = string(iter.Key())
			_, attrs = utils.SplitCompositeKey(key)
		)

		if len(attrs) < 1 {
			core.Logger.Warningf("Invalid composite key '%s': event token missing", key)
			continue
		}

		token := attrs[0]

		if err := json.Unmarshal(iter.Value(), &ticket); err != nil {
			core.Logger.Error(errors.Wrapf(err, "failed to decerialize ticket with key '%s'", key))
			continue
		}

		if now.After(ticket.Expiry) {
			dropTicketBackup(token)
			continue
		}

		rc.socketTickets[token] = ticket
	}

	if len(rc.socketTickets) != 0 {
		core.Logger.Debugf("Recovered %d event socket subscription tickets from persistence cache",
			len(rc.socketTickets),
		)
	}
}

func backupTicket(token string, ticket EventSocketSubscriptionTicket) {
	if err := core.LevelDB.Put(
		[]byte(utils.FormCompositeKey("ticket", token)),
		[]byte(utils.MustEncode(ticket)),
		nil,
	); err != nil {
		core.Logger.Error(errors.Wrapf(err, "failed to put event ticket '%s' into LevelDB", token))
	}
}

func dropTicketBackup(token string) {
	if err := core.LevelDB.Delete(
		[]byte(utils.FormCompositeKey("ticket", token)),
		nil,
	); err != nil {
		core.Logger.Error(errors.Wrapf(err, "failed to delete event tocket from LevelDB"))
	}
}


