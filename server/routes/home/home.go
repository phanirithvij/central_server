// Package home Home routes
package home

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/phanirithvij/central_server/server/config"
	"github.com/phanirithvij/central_server/server/utils"
)

var (
	// to keep track of whether the templates are initialized or not for this route
	templatesInitDone = false
	// EndpointsRegistered to keep track of whether the endpoints are registered for this route
	EndpointsRegistered = false
)

// PkgerPrefix the prefix and the top level dir for all the assets
const (
	HomeAssetsPrefix = config.PkgerPrefix + `/server/routes/home/`
)

var usage string = fmt.Sprintf(`[Error] templates are uninitialized for the %[1]s route
call %[1]s.LoadTemplates(t *template.Template) BEFORE any endpoint registrations
	eg:

	t := template.New("")
	t, err := %[1]s.LoadTemplates(t)
	if err != nil {
		log.Fatalln(err)
	}

	router.SetHTMLTemplate(t)

`, `home`)

func init() {
	// include dirs for pkger parser to pickup
	pkger.Include("/server/routes/home/home.html")
}

// TemplateParams for this route
type TemplateParams struct {
	Title string
}

// RegisterEndPoints Registers all the /api endpoints
// Must call LoadTemplates before this if it exists
// Returns the router group so it can be also used to set routes externally
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	if !templatesInitDone {
		log.Fatalln(errors.New(usage))
	}
	home := router.Group("/home")
	{
		home.GET("/", func(c *gin.Context) {
			params := TemplateParams{Title: "Home page"}
			c.HTML(http.StatusOK, HomeAssetsPrefix+"home.html", params)
		})
		home.GET("/hello", func(c *gin.Context) {
			c.String(http.StatusOK, `strings.Join(versions, "\n")`)
		})
	}
	return home
}

// Template a wrapper of template.Template
type Template struct {
	T *template.Template
}

// LoadTemplates loads the templates used by this package
func (t Template) LoadTemplates() {
	before := len(t.T.Templates())
	_, err := utils.LoadTemplates(t.T, HomeAssetsPrefix)
	if err != nil {
		log.Fatalln(err)
	}
	after := len(t.T.Templates())
	if before < after {
		templatesInitDone = true
	}
}
