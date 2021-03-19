package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/timoth-y/iot-blockchain-contracts/models"
	"github.com/timoth-y/iot-blockchain-contracts/shared"
)

// ReadingsContract provides functions for managing an models.MetricReadings from models.Device sensors
type ReadingsContract struct {
	contractapi.Contract
}

func NewReadingsContract() *ReadingsContract {
	return &ReadingsContract{}
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

func (rc *ReadingsContract) ForAsset(ctx contractapi.TransactionContextInterface, assetID string) ([]*models.MetricReadings, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("readings", []string { shared.Hash(assetID) })
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	return rc.iterate(iterator)
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
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("readings", []string { })
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return err
	}

	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		if err = ctx.GetStub().DelState(result.Key); err != nil {
			shared.Logger.Error(err)
			continue
		}
	}
	return nil
}

func (rc *ReadingsContract) iterate(iterator shim.StateQueryIteratorInterface) ([]*models.MetricReadings, error) {
	var readings []*models.MetricReadings

	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		requirement, err := models.MetricReadings{}.Decode(result.Value); if err != nil {
			shared.Logger.Error(err)
			continue
		}
		readings = append(readings, requirement)
	}

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
