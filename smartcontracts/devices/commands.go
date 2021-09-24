package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	coreutils "github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/model/couchdb"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/utils"
)

// Command handles models.DeviceCommand execution requests for devices.
// It will also log execution status in as a models.DeviceCommandLogEntry in the blockchain ledger.
func (c *DevicesContract) Command(ctx contractapi.TransactionContextInterface, payload string) error {
	req, err := requests.DeviceCommandRequest{}.Decode([]byte(payload)); if err != nil {
		return utils.LoggedError(err, "failed to deserialize request")
	}

	if err = req.Validate(); err != nil {
		return err
	}

	if exists, err := c.Exists(ctx, req.DeviceID); err != nil {
		return utils.LoggedError(err, "failed to verify device existence")
	} else if !exists {
		return errors.Errorf("device with id '%s' does not registered in the blockchain", req.DeviceID)
	}

	var key string
	if key, err = ctx.GetStub().CreateCompositeKey(couchdb.CommandLogEntryRecordType, []string{
		coreutils.Hash(req.DeviceID),
		coreutils.Hash(string(req.Command)),
		coreutils.Hash(req.IssuedAt.UTC().String()),
	}); err != nil {
		return utils.LoggedError(err, "failed to generate device command composite key")
	}

	if err = ctx.GetStub().PutState(key, couchdb.NewDeviceCommandLogEntry(&models.DeviceCommandLogEntry{
		DeviceID: req.DeviceID,
		Command: req.Command,
		Args: req.Args,
		Status: models.DeviceCmdProcessing,
		Timestamp: req.IssuedAt.UTC(),
	}).Encode()); err != nil {
		return utils.LoggedError(err, "failed to log device command in blockchain ledger")
	}

	if err = ctx.GetStub().SetEvent(fmt.Sprintf("devices.%s.command", req.DeviceID), requests.DeviceCommandEventPayload{
		ID: key,
		DeviceCommandRequest: *req,
	}.Encode()); err != nil {
		return utils.LoggedErrorf(err,
			"failed to emit '%s' command event for device '%s'", req.Command, req.DeviceID,
		)
	}

	return nil
}

// SubmitCommandResults updates models.DeviceCommandLogEntry in the blockchain ledger.
func (c *DevicesContract) SubmitCommandResults(ctx contractapi.TransactionContextInterface, entryID string, payload string) error {
	req, err := requests.DeviceCommandResultsSubmitRequest{}.Decode([]byte(payload)); if err != nil {
		return utils.LoggedError(err, "failed to deserialize request")
	}

	// Verify command source
	if data, err := ctx.GetStub().GetState(entryID); err != nil {
		return utils.LoggedError(err, "failed to verify command log entry existence")
	} else if data == nil {
		return errors.Errorf("command with id '%s' wasn't issued by devices contract, it is invalid", entryID)
	}

	data, err := ctx.GetStub().GetState(entryID); if err != nil {
		return utils.LoggedErrorf(err, "failed to retrieve command log entry with id '%s", entryID)
	}

	entry, err := models.DeviceCommandLogEntry{}.Decode(data); if err != nil {
		return utils.LoggedErrorf(err, "failed to deserialize command log entry with id '%s", entryID)
	}

	if err = req.Apply(entry); err != nil {
		return errors.Wrap(err, "failed to apply submit request on command log entry")
	}

	if err := ctx.GetStub().PutState(entryID, couchdb.NewDeviceCommandLogEntry(entry).Encode()); err != nil {
		return utils.LoggedError(err, "failed to update command log entry state")
	}

	return nil
}

// CommandsLog retrieves all models.DeviceCommandLogEntry from the blockchain ledger.
func (c *DevicesContract) CommandsLog(ctx contractapi.TransactionContextInterface, deviceID string) ([]*models.DeviceCommandLogEntry, error) {
	// Verify target device
	if exists, err := c.Exists(ctx, deviceID); err != nil {
		return nil, utils.LoggedError(err, "failed to verify device existence")
	} else if !exists {
		return nil, errors.Errorf("device with id '%s' does not registered in the blockchain", deviceID)
	}

	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(
		couchdb.CommandLogEntryRecordType,
		[]string{coreutils.Hash(deviceID)},
	)
	if err != nil {
		return nil, utils.LoggedError(err, "failed to read from world state")
	}

	var entries []*models.DeviceCommandLogEntry
	for iter.HasNext() {
		result, err := iter.Next(); if err != nil {
			core.Logger.Error(errors.Wrap(err, "failed to drain over command log results"))
			continue
		}

		entry, err := models.DeviceCommandLogEntry{}.Decode(result.Value); if err != nil {
			core.Logger.Error(errors.Wrap(err, "failed to deserialize command log entry"))
			continue
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
