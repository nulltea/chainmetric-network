package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// RequirementsContract provides functions for managing an Requirements
type RequirementsContract struct {
	contractapi.Contract
}

func NewRequirementsContract() *RequirementsContract {
	return &RequirementsContract{}
}


func (c *RequirementsContract) RetrieveRequirement(ctx contractapi.TransactionContextInterface, filter model.AssetFilter) (*model.Asset, error) {
	data, err := ctx.GetStub().GetState(filter.UniqueID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("the asset %s does not exist", filter.UniqueID)
	}

	var asset *model.Asset
	err = proto.Unmarshal(data, asset); if err != nil {
		return nil, err
	}
	return asset, nil
}

func (c *RequirementsContract) InsertRequirement(ctx contractapi.TransactionContextInterface, input model.AssetInput) (string, error) {
	id := xid.NewWithTime(time.Now()).String()
	asset := &model.Asset{
		UniqueID: id,
		Size: input.Size,
		Owner: input.Owner,
		Value: input.Value,
	}
	return id, c.InputAsset(ctx, asset)
}

func (c *RequirementsContract) InputRequirement(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	if len(asset.UniqueID) == 0 {
		return fmt.Errorf("the unique id must be setted up for asset")
	}
	data, err := proto.Marshal(asset); if err != nil {
		return err
	}
	return ctx.GetStub().PutState(asset.UniqueID, data)
}

func (c *RequirementsContract) RemoveAsset(ctx contractapi.TransactionContextInterface, filter model.AssetFilter) error {
	exists, err := c.AssetExists(ctx, filter);
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", filter.UniqueID)
	}
	return ctx.GetStub().DelState(filter.UniqueID)
}

func (c *RequirementsContract) TransferAsset(ctx contractapi.TransactionContextInterface, filter model.AssetFilter, owner string) error {
	asset, err := c.RetrieveAsset(ctx, filter); if err != nil {
		return err
	}

	asset.Owner = owner
	return c.InputAsset(ctx, asset)
}

func (c *RequirementsContract) AssetExists(ctx contractapi.TransactionContextInterface, filter model.AssetFilter) (bool, error) {
	data, err := ctx.GetStub().GetState(filter.UniqueID); if err != nil {
		return false, err
	}
	return data != nil, nil
}
