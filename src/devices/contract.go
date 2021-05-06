package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/requests"

	"github.com/timoth-y/chainmetric-contracts/shared"
)

// DevicesContract defines devices-managing Smart Contract.
type DevicesContract struct {
	contractapi.Contract
}

// NewDevicesContact creates new DevicesContract instance.
func NewDevicesContact() *DevicesContract {
	return &DevicesContract{}
}

// Retrieve retrieves models.Device from blockchain ledger.
func (c *DevicesContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Device, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	if data == nil {
		return nil, fmt.Errorf("the device %s does not exist", id)
	}

	return models.Device{}.Decode(data)
}

// All retrieves all models.Device records from blockchain ledger.
func (c *DevicesContract) All(ctx contractapi.TransactionContextInterface) ([]*models.Device, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("device", []string{})
	if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
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

// Register creates and registers new device in the blockchain ledger.
func (c *DevicesContract) Register(ctx contractapi.TransactionContextInterface, payload string) (string, error) {
	var (
		device = &models.Device{}
		err error
		event = "updated"
	)

	if device, err = device.Decode([]byte(payload)); err != nil {
		return "", shared.LoggedError(err, "failed to deserialize request")
	}

	if len(device.ID) == 0 {
		event = "inserted"

		if device.ID, err = generateCompositeKey(ctx, device); err != nil {
			return "", shared.LoggedError(err, "failed to generate composite key")
		}
	}

	if err = device.Validate(); err != nil {
		return "", errors.Wrap(err, "device is not valid")
	}

	if err := c.save(ctx, device, event); err != nil {
		return "", shared.LoggedError(err, "failed saving device")
	}

	return device.ID, nil
}

// Command handles models.DeviceCommand execution requests for devices.
// It will also log execution status in as a models.DeviceCommandLogEntry in the blockchain ledger.
func (c *DevicesContract) Command(ctx contractapi.TransactionContextInterface, payload string) error {
	req, err := requests.DeviceCommandRequest{}.Decode([]byte(payload)); if err != nil {
		return shared.LoggedError(err, "failed to deserialize request")
	}

	if err = req.Validate(); err != nil {
		return err
	}

	if exists, err := c.Exists(ctx, req.DeviceID); err != nil {
		return shared.LoggedError(err, "failed to verify device existence")
	} else if !exists {
		return errors.Errorf("device with id '%s' does not registered in the blockchain", req.DeviceID)
	}

	var key string
	if key, err = ctx.GetStub().CreateCompositeKey("command", []string{
		shared.Hash(req.DeviceID),
		xid.NewWithTime(time.Now()).String(),
	}); err != nil {
		return shared.LoggedError(err, "failed to generate device command composite key")
	}

	if err := ctx.GetStub().PutState(key, []byte(shared.MustEncode(models.DeviceCommandLogEntry{
		DeviceID: req.DeviceID,
		Command: req.Command,
		Args: req.Args,
		Status: models.DeviceCmdProcessing,
		Timestamp: time.Now().UTC(),
	}))); err != nil {
		return shared.LoggedError(err, "failed to log device command in blockchain ledger")
	}

	if err := ctx.GetStub().SetEvent("devices.command", requests.DeviceCommandEventPayload{
		ID: key,
		DeviceCommandRequest: *req,
	}.Encode()); err != nil {
		return shared.LoggedErrorf(err,
			"failed to emit '%s' command event for device '%s'", req.Command, req.DeviceID,
		)
	}

	return nil
}

// Update updates models.Device state in blockchain ledger with requested properties.
func (c *DevicesContract) Update(
	ctx contractapi.TransactionContextInterface,
	id string, payload string,
) (*models.Device, error) {
	if len(id) == 0 {
		return nil, errors.New("device id must be provided in order to update one")
	}

	device, err := c.Retrieve(ctx, id); if err != nil {
		return nil, err
	}

	req, err := requests.DeviceUpdateRequest{}.Decode([]byte(payload)); if err != nil {
		return nil, shared.LoggedError(err, "failed to deserialize request")
	}

	req.Update(device)

	if err = device.Validate(); err != nil {
		return nil, errors.Wrap(err, "device is not valid")
	}

	if err := c.save(ctx, device, "updated"); err != nil {
		return nil, shared.LoggedError(err, "failed to update device")
	}

	return device, nil
}

// Exists determines whether the models.Device exists in the blockchain ledger.
func (c *DevicesContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, err
	}
	return data != nil, nil
}

// Unbind removes models.Device from the blockchain ledger.
func (c *DevicesContract) Unbind(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.Exists(ctx, id); if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("the device with ID %q does not exist", id)
	}

	if err = ctx.GetStub().DelState(id); err != nil {
		return shared.LoggedError(err, "failed to remove device")
	}

	ctx.GetStub().SetEvent("devices.removed", models.Device{ID: id}.Encode())

	return nil
}

// RemoveAll removes all registered devices from the blockchain ledger.
// This method is for development use only and it must be removed when all dev phases will be completed.
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

		ctx.GetStub().SetEvent("devices.removed", models.Device{ID: result.Key}.Encode())
	}
	return nil
}

func (c *DevicesContract) save(
	ctx contractapi.TransactionContextInterface,
	device *models.Device,
	events ...string,
) error {
	if len(device.ID) == 0 {
		return errors.New("the unique id must be defined for device")
	}

	if err := ctx.GetStub().PutState(device.ID, device.Encode()); err != nil {
		return err
	}

	if len(events) != 0 {
		for _, event := range events {
			ctx.GetStub().SetEvent(fmt.Sprintf("devices.%s", event), device.Encode())
		}
	}

	return nil
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, dev *models.Device) (string, error) {
	return ctx.GetStub().CreateCompositeKey("device", []string{
		shared.Hash(dev.Hostname),
		xid.NewWithTime(time.Now()).String(),
	})
}
