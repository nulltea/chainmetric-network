package audience

// Notification defines data of the notification.
type Notification struct {
	Caption     string      `json:"caption"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}


