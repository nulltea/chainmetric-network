package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/timoth-y/chainmetric-contracts/model"
	"github.com/timoth-y/chainmetric-core/utils"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/requests"

	"github.com/timoth-y/chainmetric-contracts/model/response"
	"github.com/timoth-y/chainmetric-contracts/shared"
)

// AssetsContract implements assets-managing Smart Contract transaction handlers.
type AssetsContract struct {
	contractapi.Contract
}

// NewAssetsContact constructs new AssetsContract instance.
func NewAssetsContact() *AssetsContract {
	return &AssetsContract{}
}

// Retrieve retrieves models.Asset from blockchain ledger.
func (ac *AssetsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Asset, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	if data == nil {
		return nil, errors.Errorf("the asset with ID %q does not exist", id)
	}

	return models.Asset{}.Decode(data)
}

// All retrieves all models.Asset records from blockchain ledger.
func (ac *AssetsContract) All(ctx contractapi.TransactionContextInterface) ([]*models.Asset, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey("asset", []string {})
	if err != nil {
		return nil, shared.LoggedError(err, "failed to read from world state")
	}

	return ac.iterate(iter)
}

// QueryRaw performs rich query against blockchain ledger in search of specific models.Asset records.
func (ac *AssetsContract) QueryRaw(
	ctx contractapi.TransactionContextInterface,
	queryPayload string,
) ([]*models.Asset, error) {
	var (
		results []*models.Asset
		iter shim.StateQueryIteratorInterface
	)

	query, err := requests.AssetsQuery{}.Decode([]byte(queryPayload)); if err != nil {
		return nil, shared.LoggedError(err, "failed to deserialize input")
	}

	if iter, err = ctx.GetStub().GetQueryResult(queryPayload); err != nil {
		return nil,  shared.LoggedError(err, "failed to read from world state")
	}

	if results, err = ac.iterate(iter); err != nil {
		return nil, err
	}

	return results, nil
}

// Query performs rich query against blockchain ledger in search of specific models.Asset records.
//
// To support pagination it returns results wrapped in response.AssetsResponse,
// where `scroll_id` will contain special key for continuing from where the previous request ended.
func (ac *AssetsContract) Query(
	ctx contractapi.TransactionContextInterface,
	queryPayload string,
) (*response.AssetsResponse, error) {
	var (
		results []*models.Asset
		resp = &response.AssetsResponse{}
		iter shim.StateQueryIteratorInterface
	)

	query, err := requests.AssetsQuery{}.Decode([]byte(queryPayload)); if err != nil {
		return nil, shared.LoggedError(err, "failed to deserialize input")
	}

	if query.Limit != 0 {
		var md *peer.QueryResponseMetadata

		if iter, md, err = ctx.GetStub().GetQueryResultWithPagination(
			"",
			query.Limit,
			query.ScrollID,
		); err != nil {
			return nil, shared.LoggedError(err, "failed to read from world state")
		}

		resp.ScrollID = md.GetBookmark()
	} else {
		if iter, err = ctx.GetStub().GetQueryResult(""); err != nil {
			return nil, shared.LoggedError(err, "failed to read from world state")
		}
	}

	if results, err = ac.iterate(iter); err != nil {
		return nil, err
	}

	var (
		ids = make([]string, len(results))
		reqs []*models.Requirements
		reqsMap = make(map[string][]*models.Requirements)
	)

	for i := range results {
		ids[i] = results[i].ID
	}

	reqResp, err := shared.CrossChaincodeCall(ctx, "requirements", "ForAssets", utils.MustEncode(ids))
	if err != nil {
		return nil, shared.LoggedError(err, "failed requesting requirements for assets")
	}

	if err = json.Unmarshal(reqResp, &reqs); err != nil {
		return nil, shared.LoggedError(err, "failed decoding requirements for assets")
	}

	for _, req := range reqs {
		reqsMap[req.AssetID] = append(reqsMap[req.AssetID], req)
	}

	for _, asset := range results {
		if query.Satisfies(asset) {
			var r *models.Requirements

			if val, ok := reqsMap[asset.ID]; ok {
				r = val[0]
			}

			resp.Items = append(resp.Items, response.AssetResponseItem{
				Asset: *asset,
				Requirements: r,
			})
		}
	}

	return resp, nil
}

// Upsert inserts new models.Asset record into the blockchain ledger or updates existing one.
func (ac *AssetsContract) Upsert(ctx contractapi.TransactionContextInterface, data string) (string, error) {
	var (
		asset = &models.Asset{}
		event = "updated"
		err error
	)

	if asset, err = asset.Decode([]byte(data)); err != nil {
		return "", shared.LoggedError(err, "failed to deserialize input")
	}

	if len(asset.ID) == 0 {
		event = "inserted"

		if asset.ID, err = generateCompositeKey(ctx, asset); err != nil {
			return "", shared.LoggedError(err, "failed to generate composite key")
		}
	}

	return asset.ID, ac.save(ctx, asset, event)
}

// Transfer changes holder of the specific models.Asset.
func (ac *AssetsContract) Transfer(ctx contractapi.TransactionContextInterface, id string, holder string) error {
	asset, err := ac.Retrieve(ctx, id); if err != nil {
		return err
	}

	asset.Holder = holder

	return ac.save(ctx, asset)
}

// Exists determines whether the models.Asset exists in the blockchain ledger.
func (ac *AssetsContract) Exists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return false, shared.LoggedError(err, "failed to read from world state")
	}

	return data != nil, nil
}

// Remove removes models.Asset from the blockchain ledger.
func (ac *AssetsContract) Remove(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := ac.Exists(ctx, id); if err != nil {
		return err
	}

	if !exists {
		return errors.Errorf("the asset with ID %q does not exist", id)
	}

	if err := ctx.GetStub().DelState(id); err != nil {
		return shared.LoggedErrorf(err, "failed to remove asset record for id: %s", id)
	}

	return shared.LoggedError(
		ctx.GetStub().SetEvent("assets.removed", models.Asset{ID: id}.Encode()),
		"failed to emit event on asset remove",
	)
}

// RemoveAll removes all assets from the blockchain ledger.
// !! This method is for development use only and it must be removed when all dev phases will be completed.
func (ac *AssetsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("asset", []string {})
	if err != nil {
		return shared.LoggedError(err, "failed to read from world state")
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

func (ac *AssetsContract) save(ctx contractapi.TransactionContextInterface, asset *models.Asset, events ...string) error {
	if len(asset.ID) == 0 {
		return errors.New("the unique id must be defined for asset")
	}

	if err := ctx.GetStub().PutState(asset.ID, model.NewAssetRecord(asset).Encode()); err != nil {
		return shared.LoggedError(err, "failed to save asset record to blockchain ledger")
	}

	if len(events) != 0 {
		for _, event := range events {
			if err := ctx.GetStub().SetEvent(fmt.Sprintf("assets.%s", event), asset.Encode()); err != nil {
				shared.Logger.Error(errors.Wrapf(err , "failed to emit event assets.%s", event))
			}
		}
	}

	return nil
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, asset *models.Asset) (string, error) {
	return ctx.GetStub().CreateCompositeKey("asset", []string{
		shared.Hash(asset.SKU),
		xid.NewWithTime(time.Now()).String(),
	})
}

