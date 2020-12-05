package models

import (
	"log"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
)

// Server an organization server
type Server struct {
	URL   url.URL
	Alias string `validate:"alphanum"`
}

// NewServer returns a new empty organization
func NewServer() *Server {
	s := new(Server)
	defaults.SetDefaults(s)
	return s
}

// Validate validates the struct
func (s *Server) Validate() error {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			log.Println(err)
		}
	}
	return err
}
