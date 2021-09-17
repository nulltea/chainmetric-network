package model

type Notification struct {
	Caption     string            `json:"caption"`
	Description string            `json:"description"`
	Topic       string            `json:"type"`
	Audience    string            `json:"audience"`
	Data        interface{}       `json:"data"`
}


