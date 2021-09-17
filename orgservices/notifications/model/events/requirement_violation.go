package events

import (
	"fmt"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/audience"
	"github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"
)

const (
	RequirementsViolationTopic intention.EventTopic = "requirements_violation"
)

type RequirementsViolationEvent struct {
	intention.EventConcernBase
	Args struct{
		AssetID  string        `bson:"asset_id"`
		Metric   models.Metric `bson:"metric"`
	}
}

func (rv *RequirementsViolationEvent) Filter() string {
	return fmt.Sprintf("asset.%s.requirements.%s.violation", rv.Args.AssetID, rv.Args.Metric)
}

func (rv *RequirementsViolationEvent) NotificationFor(audience string, data []byte) (*audience.Notification, error) {
	return nil, nil
}

