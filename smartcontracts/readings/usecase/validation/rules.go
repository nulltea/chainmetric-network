package validation

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-network/smartcontracts/shared/utils"
)

var (
	requirements = make(map[string]map[models.Metric][]models.Requirement)
)

func SyncRequirements(ctx contractapi.TransactionContextInterface) error {
	var reqs []*models.Requirements

	payload, err := utils.CrossChaincodeCall(ctx, "requirements", "All")
	if err != nil {
		return utils.LoggedError(err, "failed to request requirements for validator")
	}

	if err = json.Unmarshal(payload, &reqs); err != nil {
		return utils.LoggedError(err, "failed to decode requirements")
	}

	for _, r := range reqs {
		if rm := requirements[r.AssetID]; rm == nil {
			requirements[r.AssetID] = make(map[models.Metric][]models.Requirement)
		}

		for m, mr := range r.Metrics {
			requirements[r.AssetID][m] = append(requirements[r.AssetID][m], mr)
		}
	}

	return nil
}

func Validate(ctx contractapi.TransactionContextInterface, r *models.MetricReadings) error {
	var reqs = requirements[r.AssetID]

	if len(reqs) == 0 {
		return nil
	}

	for m, mr := range reqs {
		for i := range mr {
			if r.Values[m] < mr[i].MinLimit || r.Values[m] > mr[i].MaxLimit {

			}
		}
	}

	return nil
}
