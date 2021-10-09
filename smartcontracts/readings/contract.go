package main

import (
	"fmt"

	"github.com/cnf/structhash"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-network/smartcontracts/readings/usecase/validation"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/model/couchdb"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/model/response"
	sharedutils "github.com/timoth-y/chainmetric-network/smartcontracts/shared/utils"

	"github.com/timoth-y/chainmetric-core/models"
)

// ReadingsContract provides functions for managing a models.MetricReadings from models.Device sensors
type ReadingsContract struct {
	contractapi.Contract
	socketTickets map[string]EventSocketSubscriptionTicket
}


// NewReadingsContract constructs new ReadingsContract instance.
func NewReadingsContract() *ReadingsContract {
	rc := &ReadingsContract{
		socketTickets: make(map[string]EventSocketSubscriptionTicket),
	}
	rc.recoverEventTicketsFromBackup()

	return rc
}

func (rc *ReadingsContract) Init(ctx contractapi.TransactionContextInterface) error {
	if err := validation.SyncRequirements(ctx); err != nil {
		return err
	}

	return nil
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
		qMap = map[string]interface{}{
			"record_type": couchdb.ReadingsRecordType,
			"asset_id":    assetID,
		}
	)

	iter, err := ctx.GetStub().GetQueryResult(core.BuildQuery(qMap, "timestamp", "asc"))
	if err != nil {
		return nil, sharedutils.LoggedError(err, "failed to read from world state")
	}

	readings := rc.drain(iter)

	for _, reading := range readings {
		for metric, value := range reading.Values {
			resp.Streams[metric] = append(resp.Streams[metric], response.MetricReadingsPoint{
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
			"record_type": couchdb.ReadingsRecordType,
			"asset_id":    assetID,
			fmt.Sprintf("values.%s", metricID): map[string]interface{}{
				"$exists": true,
			},
		}
	)

	iter, err := ctx.GetStub().GetQueryResult(core.BuildQuery(qMap, "timestamp", "asc"))
	if err != nil {
		return nil, sharedutils.LoggedError(err, "failed to read from world state")
	}

	readings := rc.drain(iter)

	for _, reading := range readings {
		if value, ok := reading.Values[metric]; ok {
			stream = append(stream, response.MetricReadingsPoint{
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
		return "", sharedutils.LoggedError(err, "failed to deserialize input")
	}

	if readings.ID, err = generateCompositeKey(ctx, readings); err != nil {
		return "", sharedutils.LoggedError(err, "failed to generate composite key")
	}

	go rc.sendToSocketListeners(ctx, readings)

	go func() {
		if err = validation.Validate(ctx, readings); err != nil {
			core.Logger.Error(err, "failed to validate readings")
		}
	}()

	return readings.ID, rc.save(ctx, readings)
}

// Exists determines whether the models.MetricReadings record exists in the blockchain ledger.
func (rc *ReadingsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, sharedutils.LoggedError(err, "failed to read from world state")
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

	return sharedutils.LoggedError(
		ctx.GetStub().DelState(id),
		"failed removing metric readings record",
	)
}

// RemoveAll removes all models.MetricReadings records from the blockchain ledger.
// !! This method is for development use only and it must be removed when all dev phases will be completed.
func (rc *ReadingsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(couchdb.ReadingsRecordType, []string {})
	if err != nil {
		return sharedutils.LoggedError(err, "failed to read from world state")
	}

	sharedutils.Iterate(iter, func(key string, _ []byte) error {
		if err = ctx.GetStub().DelState(key); err != nil {
			return errors.Wrap(err, "failed to remove readings record")
		}

		return nil
	})

	return nil
}

// NotifyRequirementsChange ...
func (rc *ReadingsContract) NotifyRequirementsChange(
	ctx contractapi.TransactionContextInterface,
	r *models.Requirements,
) {
	validation.SetRequirements(ctx, r)
}

func (rc *ReadingsContract) drain(iter shim.StateQueryIteratorInterface) []*models.MetricReadings {
	var readings []*models.MetricReadings

	sharedutils.Iterate(iter, func(_ string, value []byte) error {
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

	return ctx.GetStub().PutState(readings.ID, couchdb.NewMetricReadingsRecord(readings).Encode())
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, req *models.MetricReadings) (string, error) {
	return ctx.GetStub().CreateCompositeKey(couchdb.ReadingsRecordType, []string{
		utils.Hash(req.AssetID),
		utils.Hash(string(structhash.Dump(req, 1))),
	})
}
