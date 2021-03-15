package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/timoth-y/iot-blockchain-contracts/models"
	"github.com/timoth-y/iot-blockchain-contracts/models/request"
	"github.com/timoth-y/iot-blockchain-contracts/shared"
)

type AssetsContract struct {
	contractapi.Contract
}

func NewAssetsContact() *AssetsContract {
	return &AssetsContract{}
}

func (ac *AssetsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Asset, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("the asset with ID %q does not exist", id)
	}

	return models.Asset{}.Decode(data)
}

func (ac *AssetsContract) All(ctx contractapi.TransactionContextInterface) ([]*models.Asset, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("asset", []string {})
	if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return nil, err
	}

	return ac.iterate(iterator)
}

func (ac *AssetsContract) Query(ctx contractapi.TransactionContextInterface, query string) ([]*models.Asset, error) {
	q, err := request.AssetsQuery{}.Decode([]byte(query)); if err != nil {
		err = errors.Wrap(err, "failed to deserialize input")
		shared.Logger.Error(err)
		return nil, err
	}

	// TODO consider using CouchDB as a state database
	assets, err := ac.All(ctx); if err != nil {
		return nil, err
	}

	queried := make([]*models.Asset, 0)
	for _, asset := range assets {
		if q.Satisfies(asset) {
			queried = append(queried, asset)
		}
	}

	return queried, nil
}

func (ac *AssetsContract) Insert(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	var (
		asset = &models.Asset{}
		err error
	)

	if asset, err = asset.Decode([]byte(data)); err != nil {
		err = errors.Wrap(err, "failed to deserialize input")
		shared.Logger.Error(err)
		return "", err
	}

	if asset.ID, err = generateCompositeKey(ctx, asset); err != nil {
		err = errors.Wrap(err, "failed to generate composite key")
		shared.Logger.Error(err)
		return "", err
	}

	return asset.ID, ac.save(ctx, asset)
}

func (ac *AssetsContract) Transfer(ctx contractapi.TransactionContextInterface, id string, owner string) error {
	asset, err := ac.Retrieve(ctx, id); if err != nil {
		return err
	}
	asset.Holder = owner

	return ac.save(ctx, asset)
}

func (ac *AssetsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		err = errors.Wrap(err, "failed to read from world state")
		shared.Logger.Error(err)
		return false, err
	}

	return data != nil, nil
}

func (ac *AssetsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := ac.Exists(ctx, id); if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset with ID %q does not exist", id)
	}
	return ctx.GetStub().DelState(id)
}

func (ac *AssetsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("asset", []string {})
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

func (ac *AssetsContract) iterate(iterator shim.StateQueryIteratorInterface) ([]*models.Asset, error) {
	var assets []*models.Asset
	for iterator.HasNext() {
		result, err := iterator.Next(); if err != nil {
			shared.Logger.Error(err)
			continue
		}

		asset, err := models.Asset{}.Decode(result.Value); if err != nil {
			shared.Logger.Error(err)
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (ac *AssetsContract) save(ctx contractapi.TransactionContextInterface, asset *models.Asset) error {
	if len(asset.ID) == 0 {
		return fmt.Errorf("the unique id must be defined for asset")
	}

	if err := ctx.GetStub().PutState(asset.ID, asset.Encode()); err != nil {
		return err
	}

	return ctx.GetStub().SetEvent("assets.inserted", asset.Encode())
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, asset *models.Asset) (string, error) {
	return ctx.GetStub().CreateCompositeKey("asset", []string{
		shared.Hash(asset.SKU),
		xid.NewWithTime(time.Now()).String(),
	})
}
