package models

// Organization an organization
type Organization struct {
	Name         string
	Emails       []string
	OrgID        string
	Alias        string
	Capabilities []Capability
}
