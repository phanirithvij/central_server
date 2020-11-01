package models

// User a user who uses the organization servers
type User struct {
	Name     string
	Username string
	UID      string
	Password string
	Details  map[string]UserDetail
}

// GuestUser a guest user
type GuestUser struct {
	UID     string
	Details map[string]UserDetail
}

// UserDetail details about the user
type UserDetail struct {
	Feilds []string
}
