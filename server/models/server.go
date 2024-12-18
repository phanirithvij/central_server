package models

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Server an organization server
type Server struct {
	gorm.Model
	URL   string `validate:"url" json:"url" gorm:"uniqueindex:server_org_idx"`
	Nick  string `validate:"alphanum" json:"alias" gorm:"uniqueindex:server_org_idx"`
	OrgID uint
}

// NewServer returns a new empty server
func NewServer() *Server {
	s := new(Server)
	defaults.SetDefaults(s)
	return s
}

// Save saves to db
func (s *Server) Save(db *gorm.DB, c *gin.Context) error {
	if err := db.Create(s).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
			"type":  "create",
		})
		return err
	}
	if err := db.Save(s).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
			"type":  "save",
		})
		return err
	}
	return nil
}

// Validate validates the struct
func (s *Server) Validate() ([]string, error) {
	validate := validator.New()
	errx := validate.Struct(s)
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

// BeforeCreate before creating fix the conflicts for primarykey
func (s *Server) BeforeCreate(tx *gorm.DB) (err error) {
	cols := []clause.Column{}
	// prefix is email_
	// TODO get prefix from tx somehow?
	colsNames := []string{"updated_at", "nick", "url"}
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
