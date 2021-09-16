package model

type Notification struct {
	Caption  string           `json:"caption"`
	Topic     NotificationTopic `json:"type"`
	Audience string           `json:"audience"`
	Payload  interface{}      `json:"payload"`
}

type NotificationTopic string

const (
	RequirementsViolationTopic NotificationTopic = "requirements_violation"
)
