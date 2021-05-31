package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/timoth-y/chainmetric-contracts/model"
	"github.com/timoth-y/chainmetric-core/utils"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-contracts/shared"
)

// RequirementsContract implements requirements-managing Smart Contract transaction handlers.
type RequirementsContract struct {
	contractapi.Contract
}

// NewRequirementsContract constructs new RequirementsContract instance.
func NewRequirementsContract() *RequirementsContract {
	return &RequirementsContract{}
}

// Retrieve retrieves single models.Requirements record from blockchain ledger by a given `id`.
func (rc *RequirementsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Requirements, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	if data == nil {
		return nil, errors.Errorf("the requirement with ID %q does not exist", id)
	}

	return models.Requirements{}.Decode(data)
}

// All retrieves all models.Requirements records from blockchain ledger.
func (rc *RequirementsContract) All(ctx contractapi.TransactionContextInterface) ([]*models.Requirements, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(model.RequirementsRecordType, []string{})
	if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	return rc.drain(iter), nil
}

// ForAsset retrieves all models.Requirements records from blockchain ledger for specific asset by given `assetID`.
func (rc *RequirementsContract) ForAsset(
	ctx contractapi.TransactionContextInterface,
	assetID string,
) ([]*models.Requirements, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(model.RequirementsRecordType, []string{utils.Hash(assetID)})
	if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	return rc.drain(iter), nil
}

// ForAssets retrieves all models.Requirements records from blockchain ledger for multiply assets by given `assetIDs`.
func (rc *RequirementsContract) ForAssets(
	ctx contractapi.TransactionContextInterface,
	assetIDs []string,
) ([]*models.Requirements, error) {
	var (
		qMap = map[string]interface{}{
			"asset_id": map[string]interface{}{
				"$in": assetIDs,
			},
			"record_type": model.RequirementsRecordType,
		}
	)

	iter, err := ctx.GetStub().GetQueryResult(shared.BuildQuery(qMap, nil, nil))
	if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	return rc.drain(iter), nil
}

// Assign assigns models.Requirements to an asset and stores it in the blockchain ledger.
func (rc *RequirementsContract) Assign(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	var (
		requirements = &models.Requirements{}
		err error
		event = "updated"
	)

	if requirements, err = requirements.Decode([]byte(data)); err != nil {
		return "", shared.LoggedError(err, "failed to deserialize input")
	}

	if len(requirements.ID) == 0 {
		event = "inserted"

		if requirements.ID, err = generateCompositeKey(ctx, requirements); err != nil {
			return "", shared.LoggedError(err, "failed to generate composite key")
		}
	}

	if err = requirements.Validate(); err != nil {
		return "", errors.Wrap(err, "requirements are not valid")
	}

	if err = rc.save(ctx, requirements, event); err != nil {
		return "", shared.LoggedError(err, "failed saving requirements record")
	}

	return requirements.ID, nil
}

// Exists determines whether the models.Requirements exists in the blockchain ledger.
func (rc *RequirementsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, shared.LoggedError(err, "failed to read from world state")
	}

	return data != nil, nil
}

// Revoke revokes assignment of the models.Requirements from asset and removes it from the blockchain ledger.
func (rc *RequirementsContract) Revoke(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := rc.Exists(ctx, id); if err != nil {
		return err
	}

	if !exists {
		return errors.Errorf("the requirements with ID %q does not exist", id)
	}

	if err = ctx.GetStub().DelState(id); err != nil {
		return shared.LoggedErrorf(err, "failed to remove requirements record with id: %s", id)
	}

	if err = ctx.GetStub().SetEvent("requirements.removed", models.Requirements{ID: id}.Encode()); err != nil {
		return shared.LoggedErrorf(err, "failed to emit event on requirements remove")
	}

	return nil
}

// RemoveAll removes all registered models.Requirements records from the blockchain ledger.
// !! This method is for development use only and it must be removed when all dev phases will be completed.
func (rc *RequirementsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(model.RequirementsRecordType, []string{})
	if err != nil {
		return shared.LoggedError(err, "failed to read from world state")
	}

	shared.Iterate(iter, func(key string, _ []byte) error {
		if err = ctx.GetStub().DelState(key); err != nil {
			return errors.Wrap(err, "failed to remove requirements record")
		}

		if err = ctx.GetStub().SetEvent("requirements.removed", models.Requirements{ID: key}.Encode()); err != nil {
			return errors.Wrap(err, "failed to emit event on requirements remove")
		}

		return nil
	})

	return nil
}

func (rc *RequirementsContract) drain(iter shim.StateQueryIteratorInterface) []*models.Requirements {
	requirements := make([]*models.Requirements, 0)

	shared.Iterate(iter, func(_ string, value []byte) error {
		requirement, err := models.Requirements{}.Decode(value); if err != nil {
			return errors.Wrap(err, "failed to deserialize requirements record")
		}

		requirements = append(requirements, requirement)

		return nil
	})

	return requirements
}

func (rc *RequirementsContract) save(
	ctx contractapi.TransactionContextInterface,
	requirement *models.Requirements,
	events ...string,
) error {
	if len(requirement.ID) == 0 {
		return errors.New("the unique id must be defined for requirement")
	}

	if err := ctx.GetStub().PutState(requirement.ID, model.NewRequirementsRecord(requirement).Encode()); err != nil {
		return err
	}

	if len(events) != 0 {
		for _, event := range events {
			event = fmt.Sprintf("requirements.%s", event)
			if err := ctx.GetStub().SetEvent(
				event,
				requirement.Encode(),
			); err != nil {
				shared.Logger.Error(errors.Wrapf(err , "failed to emit event %s", event))
			}
		}
	}

	return nil
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, req *models.Requirements) (string, error) {
	return ctx.GetStub().CreateCompositeKey(model.RequirementsRecordType, []string{
		utils.Hash(req.AssetID),
		xid.NewWithTime(time.Now()).String(),
	})
}
