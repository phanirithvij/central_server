// Package models contains all the datamodels
package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	dbm "github.com/phanirithvij/central_server/server/utils/db"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = dbm.DB
)

// Organization is an organization
type Organization struct {
	gorm.Model
	OrgID string
	OrganizationPublic
	Servers []*Server `gorm:"ForeignKey:ID"`
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

// NewServer a new server for the organization
func (o *Organization) NewServer() *Server {
	s := NewServer()
	o.Servers = append(o.Servers, s)
	return s
}

// Str prints the organization
func (o *Organization) Str() string {
	jd, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintln(o)
	}
	return string(jd)
}

// SaveReq saves organization to database inside a http request
func (o *Organization) SaveReq(c *gin.Context) error {
	tx := db.Create(o)
	if tx.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": tx.Error.Error(),
			"type":  "create",
		})
		return tx.Error
	}
	tx = db.Save(o)
	if tx.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": tx.Error.Error(),
			"type":  "save",
		})
		return tx.Error
	}
	return nil
}

// Save saves
func (o *Organization) Save() error {
	tx := db.Create(o)
	if tx.Error != nil {
		return tx.Error
	}
	tx = db.Save(o)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Validate Validates the organization
func (o *Organization) Validate() ([]string, error) {
	validate := validator.New()
	errx := validate.Struct(o)
	msgs := []string{}
	if errx != nil {
		validationErrors := errx.(validator.ValidationErrors)
		for _, err := range validationErrors {
			log.Println(err, err.Field())
			msgs = append(msgs, err.Field()+" provided "+fmt.Sprint(err.Value())+" was not a valid "+strings.ToLower(err.Field()))
		}
		return msgs, errx
	}
	return []string{}, nil
}
