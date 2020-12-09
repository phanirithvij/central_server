// Package status Status routes
package status

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/phanirithvij/central_server/server/config"
	"github.com/phanirithvij/central_server/server/routes"
	"github.com/phanirithvij/central_server/server/utils"
	"github.com/phanirithvij/central_server/server/utils/sysinfo"
)

var (
	// to keep track of whether the templates are initialized or not for this route
	templatesInitDone = false
)

// PkgerPrefix the prefix and the top level dir for all the assets
const (
	statusAssets = config.PkgerPrefix + `/server/routes/` + config.Status
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

`, config.Status)

func init() {
	// include dirs for pkger parser to pickup
	pkger.Include("/server/routes/status/status.html")
}

// TemplateParams for this route
type TemplateParams struct {
	Title string
}

// SetupEndpoints Registers all the /status endpoints
// Must call LoadTemplates before this if it exists
// Returns the router group so it can be also used to set routes externally
func SetupEndpoints(router *gin.Engine) *gin.RouterGroup {
	if !templatesInitDone {
		log.Fatalln(errors.New(usage))
	}
	status := router.Group("/status")
	{
		status.GET("/", func(c *gin.Context) {
			params := TemplateParams{Title: "status page"}
			c.HTML(http.StatusOK, statusAssets+"/status.html", params)
		})
		status.GET("/json", func(c *gin.Context) {
			inf, err := sysinfo.SysInfo()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, inf)
		})
	}
	routes.RegisterSelf(config.Status)
	return status
}

// Template a wrapper of template.Template
type Template struct {
	T *template.Template
}

// LoadTemplates loads the templates used by register package
func (t Template) LoadTemplates() {
	before := len(t.T.Templates())
	_, err := utils.LoadTemplates(t.T, statusAssets)
	if err != nil {
		log.Fatalln(err)
	}
	after := len(t.T.Templates())
	if before < after {
		templatesInitDone = true
	}
}
