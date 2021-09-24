package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	coreutils "github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/core"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/model/couchdb"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/model/response"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/utils"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/requests"
)

// AssetsContract implements assets-managing Smart Contract transaction handlers.
type AssetsContract struct {
	contractapi.Contract
}

// NewAssetsContact constructs new AssetsContract instance.
func NewAssetsContact() *AssetsContract {
	return &AssetsContract{}
}

// Retrieve retrieves single models.Asset record from blockchain ledger by a given `id`.
func (ac *AssetsContract) Retrieve(ctx contractapi.TransactionContextInterface, id string) (*models.Asset, error) {
	data, err := ctx.GetStub().GetState(id); if err != nil {
		return nil, utils.LoggedError(err, "failed to read from world state")
	}

	if data == nil {
		return nil, errors.Errorf("the asset with ID %q does not exist", id)
	}

	return models.Asset{}.Decode(data)
}

// All retrieves all models.Asset records from blockchain ledger.
func (ac *AssetsContract) All(ctx contractapi.TransactionContextInterface) ([]*models.Asset, error) {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(couchdb.AssetRecordType, []string {})
	if err != nil {
		return nil, utils.LoggedError(err, "failed to read from world state")
	}

	return ac.drain(iter, nil), nil
}

// QueryRaw performs rich query against blockchain ledger in search of specific models.Asset records.
// This transaction handler will ignore 'limit' and 'scroll_id' props of the requests.AssetsQuery,
// so as not supports pagination presenting.
func (ac *AssetsContract) QueryRaw(
	ctx contractapi.TransactionContextInterface,
	queryPayload string,
) ([]*models.Asset, error) {
	query, err := requests.AssetsQuery{}.Decode([]byte(queryPayload)); if err != nil {
		return nil, utils.LoggedError(err, "failed to deserialize input")
	}

	iter, err := ctx.GetStub().GetQueryResult(buildDBQuery(query)); if err != nil {
		return nil, utils.LoggedError(err, "failed to read from world state")
	}

	return ac.drain(iter, func(a *models.Asset) bool {
		// Since location query assertion cannot be performed by state db,
		// the additional check must be performed
		return query.Location == nil || query.Location.Satisfies(a.Location)
	}), nil
}

// Query performs rich query against blockchain ledger in search of specific models.Asset records.
//
// To support pagination it returns results wrapped in response.AssetsResponse,
// where 'scroll_id' will contain special key for continuing from where the previous request ended.
func (ac *AssetsContract) Query(
	ctx contractapi.TransactionContextInterface,
	queryPayload string,
) (*response.AssetsResponse, error) {
	var (
		resp = &response.AssetsResponse{}
		iter shim.StateQueryIteratorInterface
	)

	query, err := requests.AssetsQuery{}.Decode([]byte(queryPayload)); if err != nil {
		return nil, utils.LoggedError(err, "failed to deserialize input")
	}

	if query.Limit > 0 {
		var md *peer.QueryResponseMetadata

		if iter, md, err = ctx.GetStub().GetQueryResultWithPagination(
			buildDBQuery(query),
			query.Limit,
			query.ScrollID,
		); err != nil {
			return nil, utils.LoggedError(err, "failed to read from world state")
		}

		if md.GetFetchedRecordsCount() == query.Limit {
			resp.ScrollID = md.GetBookmark()
		}
	} else {
		if iter, err = ctx.GetStub().GetQueryResult(buildDBQuery(query)); err != nil {
			return nil, utils.LoggedError(err, "failed to read from world state")
		}
	}

	results := ac.drain(iter, func(a *models.Asset) bool {
		// Since location query assertion cannot be performed by state db,
		// the additional check must be performed
		return query.Location == nil || query.Location.Satisfies(a.Location)
	})

	var (
		ids = make([]string, len(results))
		reqs []*models.Requirements
		reqsMap = make(map[string][]*models.Requirements)
	)

	for i := range results {
		ids[i] = results[i].ID
	}

	reqResp, err := utils.CrossChaincodeCall(ctx, "requirements", "ForAssets", coreutils.MustEncode(ids))
	if err != nil {
		return nil, utils.LoggedError(err, "failed requesting requirements for assets")
	}

	if err = json.Unmarshal(reqResp, &reqs); err != nil {
		return nil, utils.LoggedError(err, "failed decoding requirements for assets")
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
		return "", utils.LoggedError(err, "failed to deserialize input")
	}

	if len(asset.ID) == 0 {
		event = "inserted"

		if asset.ID, err = generateCompositeKey(ctx, asset); err != nil {
			return "", utils.LoggedError(err, "failed to generate composite key")
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
		return false, utils.LoggedError(err, "failed to read from world state")
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
		return utils.LoggedErrorf(err, "failed to remove asset record for id: %s", id)
	}

	return utils.LoggedError(
		ctx.GetStub().SetEvent("assets.removed", models.Asset{ID: id}.Encode()),
		"failed to emit event on asset remove",
	)
}

// RemoveAll removes all assets from the blockchain ledger.
// !! This method is for development use only and it must be removed when all dev phases will be completed.
func (ac *AssetsContract) RemoveAll(ctx contractapi.TransactionContextInterface) error {
	iter, err := ctx.GetStub().GetStateByPartialCompositeKey(couchdb.AssetRecordType, []string {})
	if err != nil {
		return utils.LoggedError(err, "failed to read from world state")
	}

	utils.Iterate(iter, func(key string, _ []byte) error {
		if err = ctx.GetStub().DelState(key); err != nil {
			return errors.Wrap(err, "failed to remove asset record")
		}

		if err = ctx.GetStub().SetEvent("assets.removed", models.Asset{ID: key}.Encode()); err != nil {
			return errors.Wrap(err, "failed to emit event on asset remove")
		}

		return nil
	})

	return nil
}

func (ac *AssetsContract) drain(
	iter shim.StateQueryIteratorInterface,
	predicate func(a *models.Asset) bool,
) []*models.Asset {
	var assets []*models.Asset

	utils.Iterate(iter, func(_ string, value []byte) error {
		asset, err := models.Asset{}.Decode(value); if err != nil {
			return errors.Wrap(err, "failed to deserialize asset record")
		}

		if predicate == nil || predicate(asset) {
			assets = append(assets, asset)
		}

		return nil
	})

	return assets
}

func (ac *AssetsContract) save(ctx contractapi.TransactionContextInterface, asset *models.Asset, events ...string) error {
	if len(asset.ID) == 0 {
		return errors.New("the unique id must be defined for asset")
	}

	if err := ctx.GetStub().PutState(asset.ID, couchdb.NewAssetRecord(asset).Encode()); err != nil {
		return utils.LoggedError(err, "failed to save asset record to blockchain ledger")
	}

	if len(events) != 0 {
		for _, event := range events {
			event = fmt.Sprintf("assets.%s", event)
			if err := ctx.GetStub().SetEvent(event, asset.Encode()); err != nil {
				core.Logger.Error(errors.Wrapf(err , "failed to emit event %s", event))
			}
		}
	}

	return nil
}

func buildDBQuery(req *requests.AssetsQuery) string {
	var (
		qMap = make(map[string]interface{})
	)

	if len(req.IDs) > 1 {
		qMap["asset_id"] = map[string]interface{}{
			"$in": req.IDs,
		}
	} else if len(req.IDs) == 1 {
		qMap["asset_id"] = req.IDs[0]
	}

	if req.Type != nil {
		qMap["type"] = *req.Type
	}

	if req.Holder != nil {
		qMap["holder"] = *req.Holder
	}

	if req.State != nil {
		qMap["state"] = *req.State
	}

	if len(req.Tags) > 0 {
		qMap["tags"] = map[string]interface{}{
			"$elemMatch": map[string]interface{}{
				"$in": req.Tags,
			},
		}
	}

	qMap["record_type"] = couchdb.AssetRecordType

	return core.BuildQuery(qMap)
}

func generateCompositeKey(ctx contractapi.TransactionContextInterface, asset *models.Asset) (string, error) {
	return ctx.GetStub().CreateCompositeKey(couchdb.AssetRecordType, []string{
		coreutils.Hash(asset.SKU),
		coreutils.Hash(asset.Type),
		coreutils.Hash(asset.Holder),
	})
}

