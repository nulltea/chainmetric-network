package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/timoth-y/chainmetric-contracts/model"
	"github.com/timoth-y/chainmetric-core/utils"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-contracts/model/response"
	"github.com/timoth-y/chainmetric-contracts/shared"
)

// ReadingsContract provides functions for managing an models.MetricReadings from models.Device sensors
type ReadingsContract struct {
	contractapi.Contract
	socketTickets map[string]EventSocketSubscriptionTicket
}


// NewReadingsContract constructs new ReadingsContract instance.
func NewReadingsContract() *ReadingsContract {
	return &ReadingsContract{
		socketTickets: make(map[string]EventSocketSubscriptionTicket),
	}
}

// ForAsset retrieves all models.MetricReadings records from blockchain ledger for specific asset by given `assetID`,
// aggregating them into the response.MetricReadingsStream.
func (rc *ReadingsContract) ForAsset(
	ctx contractapi.TransactionContextInterface,
	assetID string,
) (*response.MetricReadingsResponse, error) {
	var (
		resp = &response.MetricReadingsResponse{
			AssetID: assetID,
			Streams: map[models.Metric]response.MetricReadingsStream{},
		}
	)

	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(model.ReadingsRecordType, []string{utils.Hash(assetID)})
	if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	readings := rc.drain(iter)

	for _, reading := range readings {
		for metric, value := range reading.Values {
			resp.Streams[metric] = append(resp.Streams[metric], response.MetricReadingsPoint {
				DeviceID: reading.DeviceID,
				Location: reading.Location,
				Timestamp: reading.Timestamp,
				Value: value,
			})
		}
	}

	return resp, nil
}

// ForMetric retrieves all models.MetricReadings records from blockchain ledger for specific asset by given `assetID`,
// aggregating them into the response.MetricReadingsStream containing only values for given `metricID`.
func (rc *ReadingsContract) ForMetric(ctx contractapi.TransactionContextInterface, assetID string, metricID string) (response.MetricReadingsStream, error) {
	var (
		stream = response.MetricReadingsStream{}
		metric = models.Metric(metricID)
		qMap = map[string]interface{}{
			"asset_id": assetID,
			fmt.Sprintf("values.%s", metricID): map[string]interface{}{
				"$exists": true,
			},
			"record_type": model.ReadingsRecordType,
		}
	)

	iter, err := ctx.GetStub().GetQueryResult(shared.BuildQuery(qMap, nil, nil))
	if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	readings := rc.drain(iter)

	for _, reading := range readings {
		if value, ok := reading.Values[metric]; ok {
			stream = append(stream, response.MetricReadingsPoint {
				DeviceID: reading.DeviceID,
				Location: reading.Location,
				Timestamp: reading.Timestamp,
				Value: value,
			})
		}
	}

	return stream, nil
}

// Post inserts new models.MetricReadings record into the blockchain ledger.
func (rc *ReadingsContract) Post(ctx contractapi.TransactionContextInterface, payload string) (string, error) {
	var (
		readings = &models.MetricReadings{}
		err error
	)

	if readings, err = readings.Decode([]byte(payload)); err != nil {
		return "", shared.LoggedError(err, "failed to deserialize input")
	}

	if readings.ID, err = generateCompositeKey(ctx, readings); err != nil {
		return "", shared.LoggedError(err, "failed to generate composite key")
	}

	// Emitting requested events
	go func() {
		var (
			now = time.Now()
		)

		for token, request := range rc.socketTickets {
			if now.After(request.expiry) {
				delete(rc.socketTickets, token)
				shared.Logger.Debug(fmt.Sprintf("event emitter '%s' expired, currently registered: %d",
						token, len(rc.socketTickets),
					),
				)

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
	}()

	return readings.ID, rc.save(ctx, readings)
}

// Exists determines whether the models.MetricReadings record exists in the blockchain ledger.
func (rc *ReadingsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, shared.LoggedError(err, "failed to read from world state")
	}

	return data != nil, nil
}

// Remove removes models.MetricReadings record from the blockchain ledger.
func (rc *ReadingsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := rc.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("the readings with ID %q does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

// RemoveAll removes all models.MetricReadings records from the blockchain ledger.
// !! This method is for development use only and it must be removed when all dev phases will be completed.
func (rc *ReadingsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(model.ReadingsRecordType, []string { })
	if err != nil {
		return shared.LoggedError(err, "failed to read from world state")
	}

	shared.Iterate(iter, func(key string, _ []byte) error {
		if err = ctx.GetStub().DelState(key); err != nil {
			return errors.Wrap(err, "failed to remove readings record")
		}

		return nil
	})

	return nil
}

func (rc *ReadingsContract) RequestEventEmittingFor(ctx contractapi.TransactionContextInterface, assetID, metric string) string {
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

	shared.Logger.Debug(fmt.Sprintf("event emitter '%s' added, currently registered: %d", eventToken, len(rc.socketTickets)))

	return eventToken
}

func (rc *ReadingsContract) CancelEventEmitting(ctx contractapi.TransactionContextInterface, eventToken string) {
	delete(rc.socketTickets, eventToken)

	shared.Logger.Debug(fmt.Sprintf("event emitter '%s' canceled, currently registered: %d", eventToken, len(rc.socketTickets)))
}

func (rc *ReadingsContract) drain(iter shim.StateQueryIteratorInterface) []*models.MetricReadings {
	var readings []*models.MetricReadings

	shared.Iterate(iter, func(_ string, value []byte) error {
		record, err := models.MetricReadings{}.Decode(value); if err != nil {
			return errors.Wrap(err, "failed to deserialize readings record")
		}

		readings = append(readings, record)

		return nil
	})

	return readings
}

func (rc *ReadingsContract) save(ctx contractapi.TransactionContextInterface, readings *models.MetricReadings) error {
	if len(readings.ID) == 0 {
		return errors.New("the unique id must be defined for readings")
	}

	return ctx.GetStub().PutState(readings.ID, readings.Encode())
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, req *models.MetricReadings) (string, error) {
	return ctx.GetStub().CreateCompositeKey(model.ReadingsRecordType, []string{
		utils.Hash(req.AssetID),
		xid.NewWithTime(time.Now()).String(),
	})
}
