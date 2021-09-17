package presenter

import (
	"fmt"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/events"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
)

// ToEventConcerns casts SubscriptionRequest to one or multiple event concerns.
func (r *SubscriptionRequest) ToEventConcerns() ([]intention.EventConcern, error) {
	var concerns []intention.EventConcern

	switch args := r.Args.(type) {
	case *SubscriptionRequest_RequirementsViolation:
		for _, metric := range args.RequirementsViolation.Metrics {
			concerns = append(concerns, &events.RequirementsViolationEvent{
				Args: events.RequirementsViolationArgs{
					AssetID: args.RequirementsViolation.AssetID,
					Metric: models.Metric(metric),
				},
			})
		}
	default:
		return nil, fmt.Errorf("noop event is no accepted for production")
	}

	return concerns, nil
}
