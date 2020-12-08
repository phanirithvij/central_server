// Package models contains all the datamodels
package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// ValidateSub Validates only some fields for the organization
func (o *Organization) ValidateSub(only []string) ([]string, error) {
	validate := validator.New()
	errx := validate.Struct(o)
	msgs := []string{}
	onlyStr := strings.Join(only, ",")
	if errx != nil {
		validationErrors := errx.(validator.ValidationErrors)
		for _, err := range validationErrors {
			if strings.Contains(onlyStr, err.Field()) {
				log.Println(err, err.Field())
				msgs = append(msgs, err.Field()+" provided "+fmt.Sprint(err.Value())+" was not a valid "+strings.ToLower(err.Field()))
				errx = err
			} else {
				// skip validate for this field so no errors
				errx = nil
			}
		}
		return msgs, errx
	}
	return []string{}, nil
}

// OrgSubmission a submission from the clients
type OrgSubmission struct {
	ID          uint      `json:"id"`
	Address     string    `json:"address"`
	Alias       string    `json:"alias"`
	Description string    `json:"description"`
	Emails      []emailD  `json:"emails"`
	Location    []float64 `json:"location"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
}

type emailD struct {
	Email   string `json:"email"`
	Private bool   `default:"false" json:"private"`
}

// Find finds the org from db
func (s *OrgSubmission) Find() (*Organization, error) {
	o := s.Org()
	o.ID = s.ID
	tx := db.Find(&o)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return o, nil
}

// Org struct conversion, use Find() if needed from db
func (s *OrgSubmission) Org() *Organization {
	// TODO get from DB
	o := NewOrganization()
	o.Alias = s.Alias
	o.Emails = []Email{}
	for _, e := range s.Emails {
		o.Emails = append(o.Emails, Email{Email: e.Email, Private: e.Private})
	}
	o.Name = s.Name
	o.OrgDetails.LocationStr = s.Address
	if len(s.Location) == 2 {
		o.OrgDetails.LocationLL.Latitude = strconv.FormatFloat(s.Location[0], 'f', -1, 64)
		o.OrgDetails.LocationLL.Longitude = strconv.FormatFloat(s.Location[1], 'f', -1, 64)
	}
	o.OrgDetails.Description = s.Description
	return o
}

// OrgSubmission a submission for the clients
func (o *Organization) OrgSubmission() *OrgSubmission {
	s := new(OrgSubmission)
	s.Alias = o.Alias
	s.Address = o.OrgDetails.LocationStr
	s.Emails = []emailD{}
	for _, e := range o.Emails {
		x := new(emailD)
		x.Email = e.Email
		x.Private = e.Private
		s.Emails = append(s.Emails, *x)
	}
	s.Name = o.Name
	s.Address = o.OrgDetails.LocationStr
	if o.OrgDetails.LocationLL.Longitude != "" && o.OrgDetails.LocationLL.Latitude != "" {
		s.Location[0], _ = strconv.ParseFloat(o.OrgDetails.LocationLL.Latitude, 64)
		s.Location[1], _ = strconv.ParseFloat(o.OrgDetails.LocationLL.Longitude, 64)
	}
	s.Description = o.OrgDetails.Description
	return s
}
