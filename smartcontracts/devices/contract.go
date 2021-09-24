package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	coreutils "github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/model/couchdb"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/utils"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/requests"
)

// DevicesContract implements devices-managing Smart Contract transaction handlers.
type DevicesContract struct {
	contractapi.Contract
}

// NewDevicesContact creates new DevicesContract instance.
func NewDevicesContact() *DevicesContract {
	return &DevicesContract{}
}

// Retrieve retrieves single models.Device record from blockchain ledger  by a given `id`.
func (c *DevicesContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Device, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return nil, utils.LoggedError(err, "failed to read from world state")
	}

	if data == nil {
		return nil, fmt.Errorf("the device %s does not exist", id)
	}

	return models.Device{}.Decode(data)
}

// All retrieves all models.Device records from blockchain ledger.
func (c *DevicesContract) All(ctx contractapi.TransactionContextInterface) ([]*models.Device, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(couchdb.DeviceRecordType, []string{})
	if err != nil {
		return nil, utils.LoggedError(err, "failed to read from world state")
	}

	return c.drain(iter, nil), nil
}

// Register creates and registers new device in the blockchain ledger.
func (c *DevicesContract) Register(ctx contractapi.TransactionContextInterface, payload string) (string, error) {
	var (
		device = &models.Device{}
		err error
		event = "updated"
	)

	if device, err = device.Decode([]byte(payload)); err != nil {
		return "", utils.LoggedError(err, "failed to deserialize request")
	}

	if len(device.ID) == 0 {
		event = "inserted"

		if device.ID, err = generateCompositeKey(ctx, device); err != nil {
			return "", utils.LoggedError(err, "failed to generate composite key")
		}
	}

	if err = device.Validate(); err != nil {
		return "", errors.Wrap(err, "device is not valid")
	}

	if err := c.save(ctx, device, event); err != nil {
		return "", utils.LoggedError(err, "failed saving device")
	}

	return device.ID, nil
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
		return nil, utils.LoggedError(err, "failed to deserialize request")
	}

	req.Update(device)

	if err = device.Validate(); err != nil {
		return nil, errors.Wrap(err, "device is not valid")
	}

	if err := c.save(ctx, device, "updated"); err != nil {
		return nil, utils.LoggedError(err, "failed to save device record")
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
		return utils.LoggedErrorf(err, "failed to unbind device with id: %s", id)
	}

	return utils.LoggedError(
		ctx.GetStub().SetEvent("devices.removed", models.Device{ID: id}.Encode()),
		"failed to emit event on device remove",
	)
}

// RemoveAll removes all registered devices from the blockchain ledger.
// !! This method is for development use only and it must be removed when all dev phases will be completed.
func (c *DevicesContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(couchdb.DeviceRecordType, []string{})
	if err != nil {
		return utils.LoggedError(err, "failed to read from world state")
	}

	utils.Iterate(iter, func(key string, _ []byte) error {
		if err = ctx.GetStub().DelState(key); err != nil {
			return errors.Wrap(err, "failed to remove device record")
		}

		if err = ctx.GetStub().SetEvent("devices.removed", models.Device{ID: key}.Encode()); err != nil {
			return errors.Wrap(err, "failed to emit event on device remove")
		}

		return nil
	})

	return nil
}

func (c *DevicesContract) drain(
	iter shim.StateQueryIteratorInterface,
	predicate func(d *models.Device) bool,
) []*models.Device {
	var devices []*models.Device

	utils.Iterate(iter, func(_ string, value []byte) error {
		device, err := models.Device{}.Decode(value); if err != nil {
			return errors.Wrap(err, "failed to deserialize device record")
		}

		if predicate == nil || predicate(device) {
			devices = append(devices, device)
		}

		return nil
	})

	return devices
}

func (c *DevicesContract) save(
	ctx contractapi.TransactionContextInterface,
	device *models.Device,
	events ...string,
) error {
	if len(device.ID) == 0 {
		return errors.New("the unique id must be defined for device")
	}

	if err := ctx.GetStub().PutState(device.ID, couchdb.NewDeviceRecord(device).Encode()); err != nil {
		return err
	}

	if len(events) != 0 {
		for _, event := range events {
			event = fmt.Sprintf("devices.%s", event)
			if err := ctx.GetStub().SetEvent(event, device.Encode()); err != nil {
				core.Logger.Error(errors.Wrapf(err , "failed to emit event %s", event))
			}
		}
	}

	return nil
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, dev *models.Device) (string, error) {
	return ctx.GetStub().CreateCompositeKey(couchdb.DeviceRecordType, []string{
		coreutils.Hash(dev.Hostname),
		coreutils.Hash(dev.Holder),
	})
}
