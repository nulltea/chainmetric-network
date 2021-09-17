package audience

type Notification struct {
	Caption     string            `json:"caption"`
	Description string            `json:"description"`
	Topic       string            `json:"type"`
	Data        interface{}       `json:"data"`
}


