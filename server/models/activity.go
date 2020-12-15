package models

import "gorm.io/gorm"

// Activity is an activity that's stored in the DB
type Activity struct {
	gorm.Model
	Message string
	Type
}

// Type is event type
type Type string
