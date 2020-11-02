// Package models contains all the datamodels
package models

// Organization is an organization
type Organization struct {
	OrgID        string
	Capabilities []Capability
	OrganizationPublic
}

// OrganizationPublic all the public feilds that can be configured by the organization
type OrganizationPublic struct {
	Name   string   `type:"required"`
	Emails []string `type:"required"`
	// A slug
	Alias string `type:"optional"`
	OrgDetails
}

// OrgDetails the details of the organization
type OrgDetails struct {
	Location string `type:"optional"`
}
