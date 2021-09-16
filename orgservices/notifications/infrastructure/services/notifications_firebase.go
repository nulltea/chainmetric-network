package services

import (
	"firebase.google.com/go/v4"
)

// NotificationsFirebase defines service for sending notification via Firebase Cloud Messaging protocol.
type NotificationsFirebase struct {
	app *firebase.App
}

// NewNotificationsFirebase constructs new NotificationsFirebase service instance.
func NewNotificationsFirebase(app *firebase.App) *NotificationsFirebase {
	return &NotificationsFirebase{
		app: app,
	}
}

