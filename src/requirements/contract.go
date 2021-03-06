package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/timoth-y/iot-blockchain-contracts/model"
	"github.com/timoth-y/iot-blockchain-contracts/shared"
)

// RequirementsContract provides functions for managing an Requirements
type RequirementsContract struct {
	contractapi.Contract
}

func NewRequirementsContract() *RequirementsContract {
	return &RequirementsContract{}
}

func (c *RequirementsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*model.Requirements, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("the requirement with ID %q does not exist", id)
	}

	return model.Requirements{}.Decode(data)
}

func (c *RequirementsContract) ListForAsset(ctx contractapi.TransactionContextInterface, assetID string) ([]*model.Requirements, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("requirements", []string { assetID })
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	var requirements []*model.Requirements
	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		requirement, err := model.Requirements{}.Decode(result.Value); if err != nil {
			shared.Logger.Error(err)
			continue
		}
		requirements = append(requirements, requirement)
	}

	return requirements, nil
}

func (c *RequirementsContract) Insert(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	var (
		requirements = &model.Requirements{}
		err error
	)

	if err = json.Unmarshal([]byte(data), requirements); err != nil {
		err = errors.Wrap(err, "failed to deserialize input")
		shared.Logger.Error(err)
		return "", err
	}

	if requirements.ID, err = generateCompositeKey(ctx, requirements); err != nil {
		err = errors.Wrap(err, "failed to generate composite key")
		shared.Logger.Error(err)
		return "", err
	}

	return requirements.ID, c.save(ctx, requirements)
}

func (c *RequirementsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return false, err
	}

	return data != nil, nil
}

func (c *RequirementsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the requirement with ID %q does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (c *RequirementsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
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

func (c *RequirementsContract) save(ctx contractapi.TransactionContextInterface, requirement *model.Requirements) error {
	if len(requirement.ID) == 0 {
		return errors.New("the unique id must be defined for requirement")
	}

	return ctx.GetStub().PutState(requirement.ID, requirement.Encode())
}


func generateCompositeKey(ctx contractapi.TransactionContextInterface, req *model.Requirements) (string, error) {
	return ctx.GetStub().CreateCompositeKey("requirements", []string{
		req.AssetID,
		xid.NewWithTime(time.Now()).String(),
	})
}
