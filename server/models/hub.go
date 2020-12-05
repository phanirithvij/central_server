package models

import (
	"time"
)

// Hub the hub
type Hub struct {
	Admins               []Admin
	SuspensionDefaultDur time.Duration
}
