package validation

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"
)

type (
	ViolationOccurrence struct {
		Value     float64   `json:"value"`
		Timestamp time.Time `json:"timestamp"`
		Location  string    `json:"location"`
	}

	Notification struct {
		AssetID      string                `json:"asset_id"`
		Metric       models.Metric         `json:"metric"`
		Occurrences  []ViolationOccurrence `json:"occurrences"`
	}
)

var violations = make(map[string]map[models.Metric][]ViolationOccurrence)

func Validate(ctx contractapi.TransactionContextInterface, r *models.MetricReadings) error {
	if rqmCache == nil {
		SyncRequirements(ctx)
	}

	var reqs = rqmCache[r.AssetID]

	if len(reqs) == 0 {
		return nil
	}

	for m, mr := range reqs {
		for i := range mr {
			if v, ok := r.Values[m]; ok {
				if v < mr[i].MinLimit || v > mr[i].MaxLimit {
					if vm := violations[r.AssetID]; vm == nil {
						violations[r.AssetID] = make(map[models.Metric][]ViolationOccurrence)
					}

					violations[r.AssetID][m] = append(violations[r.AssetID][m], ViolationOccurrence{
						Value:     v,
						Timestamp: r.Timestamp.Round(time.Second),
						Location:  r.Location,
					})

					if vs := violations[r.AssetID][m]; len(vs) % 5 == 0 {
						if err := ctx.GetStub().SetEvent(fmt.Sprintf("asset.%s.requirements.%s.violation", r.AssetID, m),
							[]byte(utils.MustEncode(Notification{
								AssetID:     r.AssetID,
								Metric:      m,
							})),
						); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
