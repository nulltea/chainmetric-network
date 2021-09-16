package model

type Notification struct {
	Caption     string            `json:"caption"`
	Description string            `json:"description"`
	Topic       NotificationTopic `json:"type"`
	Audience    string            `json:"audience"`
	Data        interface{}       `json:"data"`
}

type NotificationTopic string

const (
	RequirementsViolationTopic NotificationTopic = "requirements_violation"
)
