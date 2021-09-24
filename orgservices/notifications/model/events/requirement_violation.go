package events

import (
	"fmt"
	"strings"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/audience"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
	"github.com/timoth-y/chainmetric-network/orgservices/shared/core"
)

const (
	RequirementsViolationEvent intention.EventKind = "requirements_violation"
)

type (
	// RequirementsViolationEventConcern implements intention.EventConcern for requirements violations events on chain.
	RequirementsViolationEventConcern struct {
		intention.EventConcernBase `bson:",inline"`
		Args RequirementsViolationArgs `bson:"args"`
	}

	// RequirementsViolationArgs defines arguments for filtering requirements violations events on chain.
	RequirementsViolationArgs struct {
		AssetID string        `bson:"asset_id"`
		Metric  models.Metric `bson:"metric,omitempty"`
		Any     bool          `bson:"any"`
	}
)

// NewRequirementsViolationConcerns constructs one or multiple RequirementsViolationEventConcern.
func NewRequirementsViolationConcerns(assetID string, metrics ...string) []intention.EventConcern {
	if len(metrics) == 0 {
		args := RequirementsViolationArgs{
			AssetID: assetID, Any: true,
		}

		return []intention.EventConcern {RequirementsViolationEventConcern{
			EventConcernBase: intention.NewEventConcernBase(
				RequirementsViolationEvent, "readings", args,
			), Args: args,
		}}
	}

	var cs []intention.EventConcern

	for i := range metrics {
		var args = RequirementsViolationArgs{
			AssetID: assetID, Metric: models.Metric(metrics[i]),
		}

		cs = append(cs, RequirementsViolationEventConcern{
			EventConcernBase: intention.NewEventConcernBase(
				RequirementsViolationEvent, "readings", args,
			), Args: args,
		})
	}

	return cs
}

// Topic ...
func (rv RequirementsViolationEventConcern) Topic() string {
	var (
		assetID = strings.ReplaceAll(rv.Args.AssetID, "\x00", "")
		metric = string(rv.Args.Metric)
	)

	if rv.Args.Any {
		return fmt.Sprintf("%s.requirements.violation", assetID)
	}

	return fmt.Sprintf("asset.%s.requirements.%s.violation", assetID, metric)
}

// Filter ...
func (rv RequirementsViolationEventConcern) Filter() string {
	if rv.Args.Any {
		return fmt.Sprintf("asset.%s.requirements.[a-z]+.violation", rv.Args.AssetID)
	}

	return fmt.Sprintf("asset.%s.requirements.%s.violation", rv.Args.AssetID, rv.Args.Metric)
}

func (rv RequirementsViolationEventConcern) NotificationWith(data []byte) (*audience.Notification, error) {
	payload, err := core.Fabric.GetContract("assets").EvaluateTransaction("Retrieve", rv.Args.AssetID)
	if err != nil {
		return nil, fmt.Errorf("faiiled to retrive asset with id '%s': %w", rv.Args.AssetID, err)
	}

	asset, err := models.Asset{}.Decode(payload)
	if err != nil {
		return nil, fmt.Errorf("faiiled to decode asset with id '%s': %w", rv.Args.AssetID, err)
	}

	return &audience.Notification{
		Caption: fmt.Sprintf("Warning: %s requiremnts viiolation", asset.SKU),
		Description: fmt.Sprintf("Latest %s readings are violating requiremnts for %s", asset.SKU),
		Data: map[string]interface{}{
			"asset_id": rv.Args.AssetID,
			"metric": rv.Args.Metric,
		},
	}, nil
}

