package audience

import "github.com/timoth-y/chainmetric-network/orgservices/notifications/model/intention"

// Notification defines data of the notification.
type Notification struct {
	Caption     string              `json:"caption"`
	Description string              `json:"description"`
	Kind        intention.EventKind `json:"event_kind"`
	Data        interface{}         `json:"data"`
}


