package presenter

import (
	"fmt"

	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/events"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
)

// ToEventConcerns casts SubscriptionRequest to one or multiple event concerns.
func (r *SubscriptionRequest) ToEventConcerns() ([]intention.EventConcern, error) {
	switch args := r.Args.(type) {
	case *SubscriptionRequest_RequirementsViolation:
		return events.NewRequirementsViolationConcerns(
			args.RequirementsViolation.AssetID,
			args.RequirementsViolation.Metrics...,
		), nil
	case *SubscriptionRequest_Noop:
		return nil, fmt.Errorf("noop event is for testing only")
	default:
		return nil, fmt.Errorf("unsupported arguments type for subscribtion")
	}
}

// NewSubscriptionResponse presents SubscriptionResponse with topics defined by given
func NewSubscriptionResponse(concerns ...intention.EventConcern) *SubscriptionResponse {
	var topics = make([]string, len(concerns))

	for i := range concerns {
		topics[i] = concerns[i].Topic()
	}

	return &SubscriptionResponse{
		Topics: topics,
	}
}
