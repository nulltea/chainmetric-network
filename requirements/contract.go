package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/rs/xid"

	"requirements/model"
)

// RequirementsContract provides functions for managing an Requirements
type RequirementsContract struct {
	contractapi.Contract
}

func NewRequirementsContract() *RequirementsContract {
	return &RequirementsContract{}
}

func (c *RequirementsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*model.Requirement, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("the requirement %s does not exist", id)
	}

	requirement := model.Requirement{}
	err = proto.Unmarshal(data, &requirement); if err != nil {
		return nil, err
	}
	return &requirement, nil
}

func (c *RequirementsContract) List(ctx contractapi.TransactionContextInterface) ([]*model.Requirement, error) {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	var requirements []*model.Requirement
	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			log.Error(err)
			continue
		}
		asset := model.Requirement{}
		err = proto.Unmarshal(result.Value, &asset); if err != nil {
			log.Error(err)
			continue
		}
		requirements = append(requirements, &asset)
	}
	return requirements, nil
}

func (c *RequirementsContract) Insert(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	requirement := &model.Requirement{}
	if err := json.Unmarshal([]byte(data), requirement); err != nil {
		return "", err
	}
	requirement.Id = xid.NewWithTime(time.Now()).String()
	return requirement.Id, c.save(ctx, requirement)
}

func (c *RequirementsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, err
	}
	return data != nil, nil
}

func (c *RequirementsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (c *RequirementsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}

	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			log.Error(err)
			continue
		}
		c.Remove(ctx, result.Key)
	}
	return nil
}

func (c *RequirementsContract) save(ctx contractapi.TransactionContextInterface, asset *model.Requirement) error {
	if len(asset.Id) == 0 {
		return fmt.Errorf("the unique id must be setted up for asset")
	}
	data, err := proto.Marshal(asset); if err != nil {
		return err
	}
	return ctx.GetStub().PutState(asset.Id, data)
}