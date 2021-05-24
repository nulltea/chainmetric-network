package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-contracts/model/response"
	"github.com/timoth-y/chainmetric-contracts/shared"
)

// ReadingsContract provides functions for managing an models.MetricReadings from models.Device sensors
type ReadingsContract struct {
	contractapi.Contract
	emitterRequests map[string]EventEmittingRequest
}

type EventEmittingRequest struct {
	assetID string
	metric models.Metric
	expiry time.Time
}

func NewReadingsContract() *ReadingsContract {
	return &ReadingsContract{
		emitterRequests: map[string]EventEmittingRequest{},
	}
}

func (rc *ReadingsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.MetricReadings, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("the readings with ID %q does not exist", id)
	}

	return models.MetricReadings{}.Decode(data)
}

func (rc *ReadingsContract) ForAsset(ctx contractapi.TransactionContextInterface, assetID string) (*response.MetricReadingsResponse, error) {
	var (
		resp = &response.MetricReadingsResponse{
			AssetID: assetID,
			Streams: map[models.Metric]response.MetricReadingsStream{},
		}
	)

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("readings", []string { shared.Hash(assetID) })
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	readings, err := rc.drain(iterator); if err != nil {
		err = errors.Wrap(err, "failed to iterate through readings data")
		shared.Logger.Error(err)
		return nil, err
	}

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


func (rc *ReadingsContract) ForMetric(ctx contractapi.TransactionContextInterface, assetID string, metricID string) (response.MetricReadingsStream, error) {
	var (
		stream = response.MetricReadingsStream{}
		metric = models.Metric(metricID)
	)

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("readings", []string { shared.Hash(assetID) })
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	readings, err := rc.drain(iterator); if err != nil {
		err = errors.Wrap(err, "failed to iterate through readings data")
		shared.Logger.Error(err)
		return nil, err
	}

	for _, reading := range readings {
		if value, ok := reading.Values[metric];  ok {
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

func (rc *ReadingsContract) Post(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	var (
		readings = &models.MetricReadings{}
		err error
	)

	if readings, err = readings.Decode([]byte(data)); err != nil {
		err = errors.Wrap(err, "failed to deserialize input")
		shared.Logger.Error(err)
		return "", err
	}

	if readings.ID, err = generateCompositeKey(ctx, readings); err != nil {
		err = errors.Wrap(err, "failed to generate composite key")
		shared.Logger.Error(err)
		return "", err
	}

	// Emitting requested events
	go func() {
		var (
			now = time.Now()
		)

		for token, request := range rc.emitterRequests {
			if now.After(request.expiry) {
				delete(rc.emitterRequests, token)
				shared.Logger.Debug(fmt.Sprintf("event emitter '%s' expired, currently registered: %d",
					token, len(rc.emitterRequests)),
				)

				continue
			}

			if request.assetID == readings.AssetID {
				if value, ok := readings.Values[request.metric];  ok {
					artifact := response.MetricReadingsPoint {
						DeviceID: readings.DeviceID,
						Location: readings.Location,
						Timestamp: readings.Timestamp,
						Value: value,
					}

					if payload, err := json.Marshal(artifact); err == nil {
						ctx.GetStub().SetEvent(token, payload)
					}
				}
			}
		}
	}()

	return readings.ID, rc.save(ctx, readings)
}

func (rc *ReadingsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return false, err
	}

	return data != nil, nil
}

func (rc *ReadingsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := rc.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the readings with ID %q does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (rc *ReadingsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey("readings", []string { })
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
		clientHash = shared.Hash(clientID)
		eventToken = fmt.Sprintf("%s.%s.%s", assetID, metric, clientHash)
		request = EventEmittingRequest{
			assetID: assetID,
			metric: models.Metric(metric),
			expiry: timestamp.Add(time.Hour * 1),
		}
	)

	rc.emitterRequests[eventToken] = request

	shared.Logger.Debug(fmt.Sprintf("event emitter '%s' added, currently registered: %d", eventToken, len(rc.emitterRequests)))

	return eventToken
}

func (rc *ReadingsContract) CancelEventEmitting(ctx contractapi.TransactionContextInterface, eventToken string) {
	delete(rc.emitterRequests, eventToken)

	shared.Logger.Debug(fmt.Sprintf("event emitter '%s' canceled, currently registered: %d", eventToken, len(rc.emitterRequests)))
}

func (rc *ReadingsContract) drain(iter shim.StateQueryIteratorInterface) ([]*models.MetricReadings, error) {
	var readings []*models.MetricReadings

	shared.Iterate(iter, func(_ string, value []byte) error {
		record, err := models.MetricReadings{}.Decode(value); if err != nil {
			return errors.Wrap(err, "failed to deserialize readings record")
		}

		readings = append(readings, record)

		return nil
	})

	return readings, nil
}

func (rc *ReadingsContract) save(ctx contractapi.TransactionContextInterface, readings *models.MetricReadings) error {
	if len(readings.ID) == 0 {
		return errors.New("the unique id must be defined for readings")
	}

	return ctx.GetStub().PutState(readings.ID, readings.Encode())
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, req *models.MetricReadings) (string, error) {
	return ctx.GetStub().CreateCompositeKey("readings", []string{
		shared.Hash(req.AssetID),
		xid.NewWithTime(time.Now()).String(),
	})
}
