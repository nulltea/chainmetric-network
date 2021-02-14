package main

import (
	"encoding/json"
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

func (c *AssetsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*model.Asset, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
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

func (c *AssetsContract) List(ctx contractapi.TransactionContextInterface) ([]*model.Asset, error) {
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

func (c *AssetsContract) Insert(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	asset := &model.Asset{}
	if err := json.Unmarshal([]byte(data), asset); err != nil {
		return "", err
	}
	asset.Id = xid.NewWithTime(time.Now()).String()
	return asset.Id, c.save(ctx, asset)
}

func (c *AssetsContract) Transfer(ctx contractapi.TransactionContextInterface, id string, owner string) error {
	asset, err := c.Retrieve(ctx, id); if err != nil {
		return err
	}
	asset.Owner = owner
	return c.save(ctx, asset)
}

func (c *AssetsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, err
	}
	return data != nil, nil
}

func (c *AssetsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := c.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (c *AssetsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
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

func (c *AssetsContract) save(ctx contractapi.TransactionContextInterface, asset *model.Asset) error {
	if len(asset.Id) == 0 {
		return fmt.Errorf("the unique id must be setted up for asset")
	}
	data, err := proto.Marshal(asset); if err != nil {
		return err
	}
	return ctx.GetStub().PutState(asset.Id, data)
}
