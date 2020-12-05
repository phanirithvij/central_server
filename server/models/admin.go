package models

import (
	"gorm.io/gorm"
)

// Admin an admin user for the hub
type Admin struct {
	gorm.Model
	Username     string
	PasswordHash string
	Name         string
	Email        string //Admin's email address
	Main         bool   //If admin is the main admin
	AddedBy      string //Username of admin who added this admin
	TimeZone     string
	Capabilites  []AdminCapabilitiy
}

// AdminCapabilitiy capability of an admin
//
// Exists to prevent all admins exporting everything
type AdminCapabilitiy struct {
	gorm.Model
	Type    string //Activity|Logs|Analytics Admin can view the page
	Allowed bool
	By      string //Username of admin who changed this
}
