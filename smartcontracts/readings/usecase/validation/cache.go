package validation

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/utils"
)

var rqmCache map[string]map[models.Metric][]models.Requirement

func SyncRequirements(ctx contractapi.TransactionContextInterface) error {
	var reqs []*models.Requirements

	payload, err := utils.CrossChaincodeCall(ctx, "requirements", "All")
	if err != nil {
		return utils.LoggedError(err, "failed to request requirements for validator")
	}

	if err = json.Unmarshal(payload, &reqs); err != nil {
		return utils.LoggedError(err, "failed to decode requirements")
	}

	rqmCache = make(map[string]map[models.Metric][]models.Requirement)

	for _, r := range reqs {
		SetRequirements(ctx, r)
	}

	return nil
}

func SetRequirements(ctx contractapi.TransactionContextInterface, r *models.Requirements) {
	if rqmCache == nil {
		SyncRequirements(ctx)

		return
	}

	if rm := rqmCache[r.AssetID]; rm == nil {
		rqmCache[r.AssetID] = make(map[models.Metric][]models.Requirement)
	}

	for m, mr := range r.Metrics {
		rqmCache[r.AssetID][m] = append(rqmCache[r.AssetID][m], mr)
	}
}
