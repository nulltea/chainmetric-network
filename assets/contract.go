package main

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/rs/xid"

	"assets/model"
)

type AssetsContract struct {
	contractapi.Contract
}

func NewAssetsContact() *AssetsContract {
	return &AssetsContract{}
}

func (c *AssetsContract) RetrieveAsset(ctx contractapi.TransactionContextInterface, filter model.AssetFilter) (*model.Asset, error) {
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

func (c *AssetsContract) InsertAsset(ctx contractapi.TransactionContextInterface, input model.AssetInput) (string, error) {
	id := xid.NewWithTime(time.Now()).String()
	asset := &model.Asset{
		UniqueID: id,
		Size: input.Size,
		Owner: input.Owner,
		Value: input.Value,
	}
	return id, c.InputAsset(ctx, asset)
}

func (c *AssetsContract) InputAsset(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	if len(asset.UniqueID) == 0 {
		return fmt.Errorf("the unique id must be setted up for asset")
	}
	data, err := proto.Marshal(asset); if err != nil {
		return err
	}
	return ctx.GetStub().PutState(asset.UniqueID, data)
}

func (c *AssetsContract) RemoveAsset(ctx contractapi.TransactionContextInterface, filter model.AssetFilter) error {
	exists, err := c.AssetExists(ctx, filter);
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", filter.UniqueID)
	}
	return ctx.GetStub().DelState(filter.UniqueID)
}

func (c *AssetsContract) TransferAsset(ctx contractapi.TransactionContextInterface, filter model.AssetFilter, owner string) error {
	asset, err := c.RetrieveAsset(ctx, filter); if err != nil {
		return err
	}

	asset.Owner = owner
	return c.InputAsset(ctx, asset)
}

func (c *AssetsContract) AssetExists(ctx contractapi.TransactionContextInterface, filter model.AssetFilter) (bool, error) {
	data, err := ctx.GetStub().GetState(filter.UniqueID); if err != nil {
		return false, err
	}
	return data != nil, nil
}
