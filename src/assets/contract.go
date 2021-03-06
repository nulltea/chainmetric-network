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

type AssetsContract struct {
	contractapi.Contract
}

func NewAssetsContact() *AssetsContract {
	return &AssetsContract{}
}

func (c *AssetsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("the asset with ID %q does not exist", id)
	}

	return model.Asset{}.Decode(data)
}

func (c *AssetsContract) List(ctx contractapi.TransactionContextInterface) ([]*model.Asset, error) {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	var assets []*model.Asset
	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		asset, err := model.Asset{}.Decode(result.Value); if err != nil {
			shared.Logger.Error(err)
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (c *AssetsContract) Insert(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	asset := &model.Asset{}
	if err := json.Unmarshal([]byte(data), asset); err != nil {
		err = errors.Wrap(err, "failed to deserialize input")
		shared.Logger.Error(err)
		return "", err
	}

	asset.ID = model.AssetID(xid.NewWithTime(time.Now()).String())

	return asset.ID.String(), c.save(ctx, asset)
}

func (c *AssetsContract) Transfer(ctx contractapi.TransactionContextInterface, id string, owner string) error {
	asset, err := c.Retrieve(ctx, id); if err != nil {
		return err
	}
	asset.Holder = owner

	return c.save(ctx, asset)
}

func (c *AssetsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return false, err
	}

	return data != nil, nil
}

func (c *AssetsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset with ID %q does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (c *AssetsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
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

func (c *AssetsContract) save(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	if len(asset.ID) == 0 {
		return fmt.Errorf("the unique id must be defined for asset")
	}

	return ctx.GetStub().PutState(asset.ID.String(), asset.Encode())
}
