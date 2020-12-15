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
	// ErrNoResultsFound no sql results were found
	ErrNoResultsFound = errors.New("no results found")
)

// Organization is an organization
type Organization struct {
	OrganizationPublic
	// bringing ID out to not prefix it
	/*
	 BUG with gorm where composite primarykey with int, text exists
	 autoIncrement is ignored for the int key
	*/
	ID uint `gorm:"primarykey;not null;autoIncrement:true;"`
	// https://gorm.io/docs/constraints.html#CHECK-Constraint
	// https://stackoverflow.com/a/5489759/8608146
	PasswordHash  string `gorm:"check:password_empty,password_hash <> '';not null;"`
	GormModelNoID `gorm:"embedded;embeddedPrefix:org_"`
	Server        *Server  `gorm:"foreignKey:OrgID"`
	DB            *gorm.DB `json:"-" gorm:"-" validate:"-"`
	// TODO support multiple servers as a backup or something
	// Servers       []*Server `gorm:"ForeignKey:ID"`
}

// OrganizationPublic all the public feilds that can be configured by the organization
type OrganizationPublic struct {
	Emails []Email `validate:"required,min=1,dive,required" gorm:"polymorphic:Organization;"`
	Name   string  `validate:"required,printascii"`
	// A slug which will be auto assigned if not chosen by them
	// Alias      string `validate:"alphanum,excludesall=!@#$%^&*()" gorm:"index;primaryKey"`
	Alias      string `validate:"alphanum,excludesall=!@#$%^&*()" gorm:"unique"`
	OrgDetails `validate:"required"`
}

// GormModelNoID capital so gorm will pickup
type GormModelNoID struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Email type it can be either public/private so we or others can contact them via email
type Email struct {
	gorm.Model       `gorm:"embedded;embeddedPrefix:email_"`
	Email            string `validate:"email" gorm:"uniqueindex:org_email_idx" json:"email"`
	OrganizationID   uint   `gorm:"uniqueindex:org_email_idx"`
	OrganizationType string `gorm:"uniqueindex:org_email_idx"`
	// https://stackoverflow.com/a/62711228/8608146
	// https://gorm.io/docs/create.html#Default-Values
	Private *bool `gorm:"not null;default:true"`
	Main    *bool `gorm:"not null;default:false"`
}

// OrgDetails the details of the organization
type OrgDetails struct {
	LocationStr string  `validate:"printascii"`
	LocationLL  LongLat `validate:"required" gorm:"embedded;embeddedPrefix:location_"`
	Description string  `validate:"required,printascii"`
	Private     *bool   `gorm:"not null;default:false"`
}

// LongLat longitude and lattitude
type LongLat struct {
	Longitude string `validate:"longitude"`
	Latitude  string `validate:"latitude"`
	Private   *bool  `gorm:"default:true"`
}

// NewOrganization returns a new empty organization
func NewOrganization() *Organization {
	o := new(Organization)
	defaults.SetDefaults(o)
	o.DB = dbm.GetDB()
	return o
}

// NewEmail returns a new empty email
func NewEmail() *Email {
	e := new(Email)
	defaults.SetDefaults(e)
	return e
}

// NewServer a new server for the organization
func (o *Organization) NewServer() *Server {
	s := NewServer()
	o.Server = s
	return s
}

// Str retuns a json representation of the organization
func (o *Organization) Str() string {
	jd, err := json.MarshalIndent(o, "", " ")
	// jd, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintln(o)
	}
	return string(jd)
}

// SaveReq saves organization to database inside a http request
func (o *Organization) SaveReq(c *gin.Context) error {
	db := o.DB
	// except alias everything else can be changed
	cols := []string{"updated_at", "name"}
	tx := db.Clauses(clause.OnConflict{
		// TODO all primaryKeys not just ID
		Columns: []clause.Column{{Name: "id"}},
		// TODO exept created_at everything
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&o)
	// remove this line after fixing the above cols logic
	log.Println("[main][WARNING]: Hardcoded feilds for Organization.Save", cols)
	if tx.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":    tx.Error.Error(),
			"type":     "create",
			"messages": []string{"Failed to create organization"},
		})
		return tx.Error
	}
	if err := db.Save(o).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":    tx.Error.Error(),
			"type":     "save",
			"messages": []string{"Failed to save to database"},
		})
		return err
	}
	return tx.Error
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
	retrr := ""
	if errx != nil {
		validationErrors := errx.(validator.ValidationErrors)
		for _, err := range validationErrors {
			if strings.Contains(onlyStr, err.Field()) {
				log.Println(err, err.Field())
				msgs = append(msgs, err.Field()+" provided "+fmt.Sprint(err.Value())+" was not a valid "+strings.ToLower(err.Field()))
				// the last error
				retrr += err.Field()
			}
		}
		if retrr != "" {
			// customising an error to make if useful for login via email or alias
			errx = errors.New(retrr)
		} else {
			// if reterr is empty => no errors in the required fields
			// so any other previous validation errors can be ignored
			errx = nil
		}
		return msgs, errx
	}
	return []string{}, nil
}

// OrgSubmissionPass a submission from the clients
type OrgSubmissionPass struct {
	OrgSubmission
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}

// OrgSubmission a submission from the clients
type OrgSubmission struct {
	ID              uint      `json:"id"`
	Address         string    `json:"address"`
	Alias           string    `json:"alias"`
	Description     string    `json:"description"`
	Emails          []EmailD  `json:"emails"`
	Location        []float64 `json:"location"`
	Name            string    `json:"name"`
	Private         *bool     `json:"private"`
	LocationPrivate *bool     `json:"privateLoc"`
	ServerURL       string    `json:"server"`
	ServerAlias     string    `json:"serverAlias"`
	ServerID        uint      `json:"serverID"`
}

// EmailD ?
type EmailD struct {
	Private *bool  `default:"false" json:"private"`
	Email   string `json:"email"`
	ID      uint   `json:"id"`
	Main    *bool  `json:"main"`
}

// PublicList a list of public organizations
func (o *Organization) PublicList() ([]*OrgSubmission, error) {
	if o.DB == nil {
		return nil, errors.New("DB was not initialized?")
	}
	orgs := []Organization{}

	tx := o.DB.Preload(clause.Associations).
		Where("private = false").Find(&orgs)

	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected <= 0 {
		return nil, ErrNoResultsFound
	}

	oss := []*OrgSubmission{}
	// TODO or don't fully load associations and query them with Where
	for _, orgx := range orgs {
		filteredEmails := []Email{}
		for _, em := range orgx.Emails {
			if em.Private != nil && !*em.Private {
				filteredEmails = append(filteredEmails, em)
			}
		}
		orgx.Emails = filteredEmails
		privLoc := orgx.LocationLL.Private
		if privLoc != nil && *privLoc {
			// if privare location remove it
			orgx.LocationLL = LongLat{}
		}
		oss = append(oss, orgx.OrgSubmission())
	}
	return oss, nil
}

// Find finds the org from db
func (o *Organization) Find() (*Organization, error) {
	if o.DB == nil {
		return nil, errors.New("DB was nil")
	}
	if o.ID == 0 {
		return nil, errors.New("ID was not set or set to zero")
	}
	// https://gorm.io/docs/preload.html#Preload-All
	tx := o.DB.Preload(clause.Associations).Find(&o)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected <= 0 {
		return nil, ErrNoResultsFound
	}
	return o, nil
}

// Find finds the org from db shortcut
func (s *OrgSubmission) Find() (*Organization, error) {
	o := s.Org()
	o.ID = s.ID
	return o.Find()
}

// BeforeCreate ..
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	// log.Println(tx.Statement.FullSaveAssociations)
	return nil
}

// BeforeUpdate hook to determine if any thing needs to be done
func (o *Organization) BeforeUpdate(tx *gorm.DB) error {
	// TODO if organization type changed add it to public clients?
	log.Println(tx.Statement.Changed("Private"))
	// TODO if main email changed send email verification
	// THis will not work because gorm is INSERTing emails on update
	// So do it in BeforeCreate of Email{}
	log.Println(tx.Statement.Changed("Emails"))
	return nil
}

// NewUpdate updates with new values
func (o *Organization) NewUpdate(n *Organization) error {
	// alias, ID are readonly so don't update
	n.Alias = o.Alias
	n.ID = o.ID
	// remove any emails
	removeList := []uint{}
	// orEmails := []string{}
	// newEmails := []string{}
	// for _, or := range o.Emails {
	// 	orEmails = append(orEmails, or.Email)
	// }
	// for _, ne := range n.Emails {
	// 	newEmails = append(newEmails, ne.Email)
	// }

	// TODO look at https://gorm.io/docs/associations.html#Append-Associations
	// It provides Append, Replace etc.

	for _, or := range o.Emails {
		del := true
		for _, ne := range n.Emails {
			if or.Email == ne.Email {
				// if email matched => found old email in new list
				// so not deleted
				del = false
				break
			}
		}
		if del {
			if or.Main != nil {
				if !*or.Main {
					// full loop executed so not found add it to deleted
					log.Println("Deleting email", or.Email)
					removeList = append(removeList, or.ID)
				}
				// else skip if main email
			} else {
				log.Println("Deleting email", or.Email)
				removeList = append(removeList, or.ID)
			}
		}
	}
	log.Println("Deleting emails with ids", removeList)

	// these are needed for some reason
	o.Emails = n.Emails
	o.Server = n.Server

	// updates and insterts will be done by gorm
	o.LocationLL = n.LocationLL

	if err := o.DB.Model(&o).
		// this feels optional
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&n).Error; err != nil {
		return err
	}

	if len(removeList) > 0 {
		// Perma delete emails because user may add it again
		// then non conflict should occur
		if err := o.DB.
			Unscoped().
			Delete(&Email{}, removeList).Error; err != nil {
			return err
		}
	}
	return nil
}

// FindByAlias finds the org from db by alias
func (s *OrgSubmission) FindByAlias() (*Organization, error) {
	o := s.Org()
	db := o.DB
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

/*
	Note:
	Because we need both the email and the organization
	It will not work by default as createdAt, updatedAt, deleteAt
	fields are the same for both the models which are embedded from gorm.Model
	So I create a custom struct and gave an embedded_prefix for both the models
	So there will be no column name conflict or duplicates as far as
	these 3 feilds are concerned.

	And I realized this was a huge waste of time to save two querys
	Because we do need to fetch all the emails and other associations
	for the setting page. But that would require only the org ID and nothing else

	So leaving this code as is in the hope that it will be useful somewhere else
	But only the ID from all this is needed
*/
// type orgEmailPacked struct {
// 	Organization
// 	Email
// }

// FindByEmail finds the org from db by email
func (s *OrgSubmission) FindByEmail() (*Organization, error) {
	db := dbm.GetDB()
	o := &Organization{}
	// https://gorm.io/docs/query.html#Joins
	// tx := db.Model(&Organization{}).
	// 	Select("`organizations`.*, `emails`.*").
	// 	Joins("LEFT JOIN `emails` ON `emails`.`organization_id` = `organizations`.`id`").
	// 	Where("`emails`.`email` = ?", s.Emails[0].Email).
	// 	Scan(&res)
	tx := db.Model(&Organization{}).
		Select("`organizations`.`id`, `organizations`.`password_hash`").
		Joins("LEFT JOIN `emails` ON `emails`.`organization_id` = `organizations`.`id`").
		Where(
			"`emails`.`email` = ? AND `emails`.`main` = ?",
			s.Emails[0].Email,
			true,
		).
		Scan(&o)
	// TODO need one more query anyway to get the full associations

	// so it should Ideally be just check if the email exists
	// then get the org_id and call it done

	if tx.Error != nil {
		return nil, tx.Error
	}
	// can also be -1 for non-existant emails
	if tx.RowsAffected <= 0 {
		return nil, ErrNoResultsFound
	}
	// o := &(res.Organization)
	// o.Emails = []Email{res.Email}
	return o, nil
}

// BeforeCreate before creating fix the conflicts for primarykey
func (b *Email) BeforeCreate(tx *gorm.DB) (err error) {
	cols := []clause.Column{}
	// prefix is email_
	// TODO get prefix from tx somehow?
	colsNames := []string{"email_updated_at", "email", "private"}
	for _, field := range tx.Statement.Schema.PrimaryFields {
		cols = append(cols, clause.Column{Name: field.DBName})
		colsNames = append(colsNames, field.DBName)
	}
	// https://gorm.io/docs/create.html#Upsert-On-Conflict
	// https://github.com/go-gorm/gorm/issues/3611#issuecomment-729673788
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   cols,
		DoUpdates: clause.AssignmentColumns(colsNames),
		// DoNothing: true,
	})
	return nil
}

// BeforeUpdate before updating fix the conflicts for primarykey
func (b *Email) BeforeUpdate(tx *gorm.DB) (err error) {
	// TODO this is not getting called for some reason idk why
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
	o := NewOrganization()
	o.ID = s.ID
	o.Alias = s.Alias
	o.Emails = []Email{}
	for _, e := range s.Emails {
		// TODO private is true always initially bug
		if e.Private != nil {
			log.Println(*e.Private)
		}
		o.Emails = append(o.Emails,
			Email{
				Email:   e.Email,
				Private: e.Private,
				Model:   gorm.Model{ID: e.ID},
				Main:    e.Main,
			},
		)
	}
	o.Name = s.Name
	o.OrgDetails.LocationStr = s.Address
	o.NewServer()
	o.Server.Nick = s.ServerAlias
	o.Server.URL = s.ServerURL
	o.Server.ID = s.ServerID
	o.OrgDetails.LocationLL.Private = s.LocationPrivate
	o.OrgDetails.Private = s.Private
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
	s.ID = o.ID
	s.Alias = o.Alias
	s.Address = o.OrgDetails.LocationStr
	s.Emails = []EmailD{}
	for _, e := range o.Emails {
		x := new(EmailD)
		x.Email = e.Email
		x.Private = e.Private
		x.ID = e.ID
		x.Main = e.Main
		s.Emails = append(s.Emails, *x)
	}
	s.Name = o.Name
	s.Address = o.OrgDetails.LocationStr

	s.LocationPrivate = o.OrgDetails.LocationLL.Private
	s.Private = o.OrgDetails.Private

	// single server info
	s.ServerAlias = o.Server.Nick
	s.ServerURL = o.Server.URL
	s.ServerID = o.Server.ID

	if o.OrgDetails.LocationLL.Longitude != "" && o.OrgDetails.LocationLL.Latitude != "" {
		s.Location = []float64{0, 0}
		s.Location[0], _ = strconv.ParseFloat(o.OrgDetails.LocationLL.Latitude, 64)
		s.Location[1], _ = strconv.ParseFloat(o.OrgDetails.LocationLL.Longitude, 64)
	}
	s.Description = o.OrgDetails.Description
	return s
}
