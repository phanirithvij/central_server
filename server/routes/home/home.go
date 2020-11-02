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
	"github.com/phanirithvij/central_server/server/utils"
)

// templatesInitDone to keep track of whether the templates are initialized or not for this route
var templatesInitDone = false

// PkgerPrefix the prefix and the top level dir for all the assets
const (
	PkgerPrefix      = "github.com/phanirithvij/central_server:"
	HomeAssetsPrefix = PkgerPrefix + `/server/routes/home/`
)

var usage string = fmt.Sprintf(`[Error] templates are uninitialized for the %[1]s route
call %[1]s.LoadTemplates(t *template.Template) BEFORE any endpoint registrations
	eg:

	t := template.New("")
	t, err := home.LoadTemplates(t)
	if err != nil {
		log.Fatalln(err)
	}

	router.SetHTMLTemplate(t)

`, `home`)

func init() {
	// include dirs for pkger parser to pickup
	pkger.Include("/server/routes/home/home.html")
}

// RegisterEndPoints Registers all the /api endpoints
// Must call LoadTemplates before this if it exists
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	if !templatesInitDone {
		log.Fatalln(errors.New(usage))
	}
	home := router.Group("/home")
	{
		home.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, HomeAssetsPrefix+"home.html", gin.H{
				"title": "Home Page",
			})
		})
		home.GET("/hello", func(c *gin.Context) {
			c.String(http.StatusOK, `strings.Join(versions, "\n")`)
		})
	}
	return home
}

// https://gin-gonic.com/docs/examples/bind-single-binary-with-template/
// https://github.com/gin-gonic/examples/commit/c5a87f03d39fdb9e0f6312344c21ccdd55140293

// Template an alias of template.Template
type Template struct {
	*template.Template
}

// LoadTemplates loads the templates used by this package
func (t *Template) LoadTemplates() {
	templatesInitDone = true
	var s interface{} = t
	_, err := utils.LoadTemplates(s.(*template.Template), HomeAssetsPrefix)
	if err != nil {
		log.Fatalln(err)
	}
}
