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

// RequirementsContract provides functions for managing an Requirements
type RequirementsContract struct {
	contractapi.Contract
}

func NewRequirementsContract() *RequirementsContract {
	return &RequirementsContract{}
}

func (rc *RequirementsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Requirements, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("the requirement with ID %q does not exist", id)
	}

	return models.Requirements{}.Decode(data)
}

func (rc *RequirementsContract) All(ctx contractapi.TransactionContextInterface) ([]*models.Requirements, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("requirements", []string{})
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	return rc.iterate(iterator)
}

func (rc *RequirementsContract) ForAsset(ctx contractapi.TransactionContextInterface, assetID string) ([]*models.Requirements, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("requirements", []string { shared.Hash(assetID) })
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	return rc.iterate(iterator)
}

func (rc *RequirementsContract) ForAssets(ctx contractapi.TransactionContextInterface, assetIDs []string) ([]*models.Requirements, error) {
	var (
		results = make([]*models.Requirements, 0)
	)

	for i := range assetIDs {
		reqs, err := rc.ForAsset(ctx, assetIDs[i]); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		results = append(results, reqs...)
	}

	return results, nil
}

func (rc *RequirementsContract) Assign(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	var (
		requirements = &models.Requirements{}
		err error
		event = "updated"
	)

	if requirements, err = requirements.Decode([]byte(data)); err != nil {
		err = errors.Wrap(err, "failed to deserialize input")
		shared.Logger.Error(err)
		return "", err
	}

	if len(requirements.ID) == 0 {
		event = "inserted"

		if requirements.ID, err = generateCompositeKey(ctx, requirements); err != nil {
			err = errors.Wrap(err, "failed to generate composite key")
			shared.Logger.Error(err)
			return "", err
		}
	}

	if err = requirements.Validate(); err != nil {
		return "", errors.Wrap(err, "requirements are not valid")
	}

	return requirements.ID, rc.save(ctx, requirements, event)
}

func (rc *RequirementsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return false, err
	}

	return data != nil, nil
}

func (rc *RequirementsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := rc.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the requirement with ID %q does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (rc *RequirementsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("requirements", []string{})
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

		ctx.GetStub().SetEvent("requirements.removed", models.Requirements{ID: result.Key}.Encode())
	}
	return nil
}

func (rc *RequirementsContract) iterate(iterator shim.StateQueryIteratorInterface) ([]*models.Requirements, error) {
	var requirements []*models.Requirements

	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		requirement, err := models.Requirements{}.Decode(result.Value); if err != nil {
			shared.Logger.Error(err)
			continue
		}
		requirements = append(requirements, requirement)
	}

	return requirements, nil
}

func (rc *RequirementsContract) save(ctx contractapi.TransactionContextInterface, requirement *models.Requirements, events ...string) error {
	if len(requirement.ID) == 0 {
		return errors.New("the unique id must be defined for requirement")
	}

	if err := ctx.GetStub().PutState(requirement.ID, requirement.Encode()); err != nil {
		return err
	}

	if len(events) != 0 {
		for _, event := range events {
			ctx.GetStub().SetEvent(fmt.Sprintf("requirements.%s", event), requirement.Encode())
		}
	}

	return nil
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, req *models.Requirements) (string, error) {
	return ctx.GetStub().CreateCompositeKey("requirements", []string{
		shared.Hash(req.AssetID),
		xid.NewWithTime(time.Now()).String(),
	})
}
