// Package config contain various configurations
package config

const (
	// PkgerPrefix ..
	PkgerPrefix = "github.com/phanirithvij/central_server:"
	Register    = "register"
	Logout      = "logout"
	Login       = "login"
	Status      = "status"
	Settings    = "settings"
	Home        = "home"
	API         = "api"
)

// EndPointStrings ...
var EndPointStrings = []string{
	API,
	Home,
	Register,
	Settings,
	Status,
	Logout,
	Login,
}
