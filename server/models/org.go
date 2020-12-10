// Package models contains all the datamodels
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	dbm "github.com/phanirithvij/central_server/server/utils/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNoResultsFound = errors.New("No results found")
)

// Organization is an organization
type Organization struct {
	gorm.Model
	OrganizationPublic
	Servers []*Server `gorm:"ForeignKey:ID"`
	DB      *gorm.DB  `json:"-" gorm:"-" validate:"-"`
}

// OrganizationPublic all the public feilds that can be configured by the organization
type OrganizationPublic struct {
	Emails []Email `validate:"required,min=1,dive,required" gorm:"polymorphic:Organization;"`
	Name   string  `validate:"required,printascii"`
	// A slug which will be auto assigned if not chosen by them
	Alias      string `validate:"alphanum" gorm:"index"`
	OrgDetails `validate:"required"`
}

// Email type it can be either public/private so we or others can contact them via email
type Email struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Email            string         `validate:"email" gorm:"uniqueindex:org_email_idx;primaryKey" json:"email"`
	OrganizationID   uint           `gorm:"uniqueindex:org_email_idx;primaryKey"`
	OrganizationType string         `gorm:"primaryKey"`
	Private          bool           `default:"true"`
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
	o.DB = dbm.GetDB()
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
	db := o.DB
	if err := db.Create(o).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":    err.Error(),
			"type":     "create",
			"messages": []string{"Failed to create organization"},
		})
		return err
	}
	if err := db.Save(o).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":    err.Error(),
			"type":     "save",
			"messages": []string{"Failed to save to database"},
		})
		return err
	}
	return nil
}

// // Save saves
// func (o *Organization) Save() error {
// 	db := o.DB
// 	if err := db.Create(o).Error; err != nil {
// 		return err
// 	}
// 	if err := db.Save(o).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

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

// OrgSubmissionPass a submission from the clients
type OrgSubmissionPass struct {
	OrgSubmission
	Password string `json:"password"`
}

// OrgSubmission a submission from the clients
type OrgSubmission struct {
	ID          uint      `json:"id"`
	Address     string    `json:"address"`
	Alias       string    `json:"alias"`
	Description string    `json:"description"`
	Emails      []EmailD  `json:"emails"`
	Location    []float64 `json:"location"`
	Name        string    `json:"name"`
}

// EmailD ?
type EmailD struct {
	Email   string `json:"email"`
	Private bool   `default:"false" json:"private"`
}

// Find finds the org from db
func (s *OrgSubmission) Find() (*Organization, error) {
	o := s.Org()
	o.ID = s.ID
	db := o.DB
	log.Println(o.Str())
	// https://gorm.io/docs/preload.html#Preload-All
	tx := db.Preload(clause.Associations).Find(&o)
	if tx.Error != nil {
		return nil, tx.Error
	}
	log.Println(tx.RowsAffected)
	return o, nil
}

// FindByAlias finds the org from db by alias
func (s *OrgSubmission) FindByAlias() (*Organization, error) {
	o := s.Org()
	db := o.DB
	log.Println(o.Str())
	// https://gorm.io/docs/preload.html#Preload-All
	tx := db.Preload(clause.Associations).Where("alias = ?", o.Alias).Find(&o)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, ErrNoResultsFound
	}
	return o, nil
}

// FindByEmail finds the org from db by email
func (s *OrgSubmission) FindByEmail() (*Organization, error) {
	o := s.Org()
	db := o.DB
	e := new(Email)
	e.Email = s.Emails[0].Email

	// https://gorm.io/docs/preload.html#Preload-All
	tx := db.Find(&e)
	if tx.Error != nil {
		return nil, tx.Error
	}
	log.Println(tx.RowsAffected)
	log.Println(o.Str())
	return o, nil
}

// Save a org
//
// Upserts the org
func (o *Organization) Save() error {
	db := o.DB
	cols := []string{"updated_at", "name"}
	// PGSQL
	// cols := []string{"updated_at", "name", "emails"}
	tx := db.Clauses(clause.OnConflict{
		// TODO all primaryKeys not just ID
		Columns: []clause.Column{{Name: "id"}},
		// TODO exept created_at everything
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&o)
	// remove this line after fixing the above cols logic
	log.Println("[main][WARNING]: Hardcoded feilds for Organization.Save", cols)
	return tx.Error
}

// BeforeCreate before creating fix the conflicts for primarykey
func (b *Email) BeforeCreate(tx *gorm.DB) (err error) {
	cols := []clause.Column{}
	colsNames := []string{}
	for _, field := range tx.Statement.Schema.PrimaryFields {
		cols = append(cols, clause.Column{Name: field.DBName})
		colsNames = append(colsNames, field.DBName)
	}
	// https://gorm.io/docs/create.html#Upsert-On-Conflict
	// https://github.com/go-gorm/gorm/issues/3611#issuecomment-729673788
	tx.Statement.AddClause(clause.OnConflict{
		Columns: cols,
		// DoUpdates: clause.AssignmentColumns(colsNames),
		DoNothing: true,
	})
	return nil
}

// BeforeUpdate before updating fix the conflicts for primarykey
func (b *Email) BeforeUpdate(tx *gorm.DB) (err error) {
	cols := []clause.Column{}
	colsNames := []string{}
	for _, field := range tx.Statement.Schema.PrimaryFields {
		cols = append(cols, clause.Column{Name: field.DBName})
		colsNames = append(colsNames, field.DBName)
	}
	colsNames = append(colsNames, "updated_at")
	// https://gorm.io/docs/create.html#Upsert-On-Conflict
	// https://github.com/go-gorm/gorm/issues/3611#issuecomment-729673788
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   cols,
		DoUpdates: clause.AssignmentColumns(colsNames),
	})
	return nil
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
	s.Emails = []EmailD{}
	for _, e := range o.Emails {
		x := new(EmailD)
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
