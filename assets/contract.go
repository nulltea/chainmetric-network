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

func (c *AssetsContract) RetrieveAsset(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	asset := model.Asset{}
	err = proto.Unmarshal(data, &asset); if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (c *AssetsContract) ListAssets(ctx contractapi.TransactionContextInterface) ([]*model.Asset, error) {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	var assets []*model.Asset
	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			log.Error(err)
			continue
		}
		asset := model.Asset{}
		err = proto.Unmarshal(result.Value, &asset); if err != nil {
			log.Error(err)
			continue
		}
		assets = append(assets, &asset)
	}
	return assets, nil
}

func (c *AssetsContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []*model.Asset{
		&model.Asset{
			UniqueID: xid.NewWithTime(time.Now()).String(),
			Name: "Asset1",
			Description: "Asset1",
			Size_: "1.0",
			Owner: "admin",
			Value: "1.0",
		},
		&model.Asset{
			UniqueID: xid.NewWithTime(time.Now()).String(),
			Name: "Asset2",
			Description: "Asset2",
			Size_: "1.5",
			Owner: "admin",
			Value: "5.0",
		},
	}
	for _, asset := range assets {
		c.inputAsset(ctx, asset)
	}
	return nil
}

func (c *AssetsContract) InsertAsset(ctx contractapi.TransactionContextInterface, name, desc, size, owner, value string) (string, error) {
	id := xid.NewWithTime(time.Now()).String()
	asset := &model.Asset{
		UniqueID: id,
		Name: name,
		Description: desc,
		Size_: size,
		Owner: owner,
		Value: value,
	}
	return id, c.inputAsset(ctx, asset)
}

func (c *AssetsContract) inputAsset(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	if len(asset.UniqueID) == 0 {
		return fmt.Errorf("the unique id must be setted up for asset")
	}
	data, err := proto.Marshal(asset); if err != nil {
		return err
	}
	return ctx.GetStub().PutState(asset.UniqueID, data)
}

func (c *AssetsContract) RemoveAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.AssetExists(ctx, id);
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (c *AssetsContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, owner string) error {
	asset, err := c.RetrieveAsset(ctx, id); if err != nil {
		return err
	}
	asset.Owner = owner
	return c.inputAsset(ctx, asset)
}

func (c *AssetsContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, err
	}
	return data != nil, nil
}
