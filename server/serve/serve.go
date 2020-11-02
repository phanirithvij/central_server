// Package serve serves the server
package serve

import (
	"html/template"
	"log"
	"strconv"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/models"
	routes "github.com/phanirithvij/central_server/server/routes"
	api "github.com/phanirithvij/central_server/server/routes/api"
	home "github.com/phanirithvij/central_server/server/routes/home"
	register "github.com/phanirithvij/central_server/server/routes/register"
	"github.com/phanirithvij/central_server/server/utils"
)

// Serve A function which serves the server
func Serve(port int, debug bool) {
	if debug {
		log.SetFlags(log.Ltime | log.Lshortfile)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	registerTemplates(router)

	api.RegisterEndPoints(router)
	home.RegisterEndPoints(router)
	register.RegisterEndPoints(router)

	routes.CheckEndpoints()

	// printStruct()

	endless.ListenAndServe(":"+strconv.Itoa(port), router)
}

func registerTemplates(router *gin.Engine) {

	t := template.New("")
	ht := home.Template{T: t}
	ht.LoadTemplates()

	rt := register.Template{T: t}
	rt.LoadTemplates()

	router.SetHTMLTemplate(t)
}

func printStruct() {

	o := models.Organization{
		OrgID:        "org-oror",
		Capabilities: []models.Capability{},
		OrganizationPublic: models.OrganizationPublic{
			Alias: "oror",
			Emails: []string{
				"hello@kk",
				"hello@kk",
				"hello@kk",
				"hello@kk",
				"hello@kk",
			},
			Name: "Or Or Organization",
			OrgDetails: models.OrgDetails{
				Location: "Hyd",
			},
		},
	}
	utils.PrintStruct(o)
}
