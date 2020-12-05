// Package models contains all the datamodels
package models

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"gorm.io/gorm"
)

// Organization is an organization
type Organization struct {
	gorm.Model
	OrgID string
	OrganizationPublic
}

// OrganizationPublic all the public feilds that can be configured by the organization
type OrganizationPublic struct {
	Name   string  `validate:"required,printascii"`
	Emails []Email `validate:"required,min=1,dive,required" gorm:"ForeignKey:ID"`
	// A slug which will be auto assigned if not chosen by them
	Alias      string `validate:"alphanum"`
	OrgDetails `validate:"required"`
}

// Email type it can be either public/private so we or others can contact them via email
type Email struct {
	gorm.Model
	Email   string `validate:"email"`
	Private bool   `default:"true"`
}

// OrgDetails the details of the organization
type OrgDetails struct {
	LocationStr string  `validate:"printascii"`
	LocationLL  LongLat `validate:"required" gorm:"embedded;embeddedPrefix:location_"`
	Description string  `validate:"required,alphanumunicode"`
	Private     bool    `default:"false"`
}

// LongLat longitude and lattitude
type LongLat struct {
	Longitude string `validate:"longitude"`
	Latitude  string `validate:"latitude"`
	Private   bool   `default:"true"`
}

// NewOrganization returns a new empty organization
func NewOrganization() *Organization {
	o := new(Organization)
	defaults.SetDefaults(o)
	return o
}

// Print prints the organization
func (o *Organization) Print() {
	jd, err := json.Marshal(o)
	if err != nil {
		fmt.Println(o)
	}
	fmt.Println(string(jd))
}

// Validate Validates the organization
func (o *Organization) Validate() error {
	validate := validator.New()
	err := validate.Struct(o)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			log.Println(err)
		}
	}
	return err
}
