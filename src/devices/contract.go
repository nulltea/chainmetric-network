package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/timoth-y/iot-blockchain-contracts/models"
	"github.com/timoth-y/iot-blockchain-contracts/shared"
)

type DevicesContract struct {
	contractapi.Contract
}

func NewDevicesContact() *DevicesContract {
	return &DevicesContract{}
}

func (c *DevicesContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Device, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("the device %s does not exist", id)
	}

	return models.Device{}.Decode(data)
}

func (c *DevicesContract) List(ctx contractapi.TransactionContextInterface) ([]*models.Device, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("device", []string{})
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	var devices []*models.Device
	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		device, err := models.Device{}.Decode(result.Value); if err != nil {
			shared.Logger.Error(err)
			continue
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (c *DevicesContract) Insert(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	var (
		device = &models.Device{}
		err error
	)

	if device, err = device.Decode([]byte(data)); err != nil {
		err = errors.Wrap(err, "failed to deserialize input")
		shared.Logger.Error(err)
		return "", err
	}

	if device.ID, err = generateCompositeKey(ctx, device); err != nil {
		err = errors.Wrap(err, "failed to generate composite key")
		shared.Logger.Error(err)
		return "", err
	}

	return device.ID, c.save(ctx, device)
}

func (c *DevicesContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, err
	}
	return data != nil, nil
}

func (c *DevicesContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the device with ID %q does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (c *DevicesContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("device", []string{})
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

func (c *DevicesContract) save(ctx contractapi.TransactionContextInterface, device *models.Device) error {
	if len(device.ID) == 0 {
		return fmt.Errorf("the unique id must be defined for device")
	}
	return ctx.GetStub().PutState(device.ID, device.Encode())
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, dev *models.Device) (string, error) {
	return ctx.GetStub().CreateCompositeKey("device", []string{
		dev.Hostname,
		xid.NewWithTime(time.Now()).String(),
	})
}
