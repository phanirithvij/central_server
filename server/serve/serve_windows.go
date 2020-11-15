// Package serve serves the server
package serve

import (
	"html/template"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/models"
	routes "github.com/phanirithvij/central_server/server/routes"
	api "github.com/phanirithvij/central_server/server/routes/api"
	home "github.com/phanirithvij/central_server/server/routes/home"
	register "github.com/phanirithvij/central_server/server/routes/register"
	status "github.com/phanirithvij/central_server/server/routes/status"
	"github.com/phanirithvij/central_server/server/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
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
	status.RegisterEndPoints(router)

	routes.CheckEndpoints()

	o := newOrg()
	utils.PrintStruct(*o)
	// o.Print()
	o.Validate()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Organization{}, models.Email{})
	if err != nil {
		log.Fatalln(err)
	}

	db.Create(&o)
	db.Save(&o)

	http.ListenAndServe(":"+strconv.Itoa(port), router)
}

func registerTemplates(router *gin.Engine) {

	t := template.New("")
	ht := home.Template{T: t}
	ht.LoadTemplates()

	rt := register.Template{T: t}
	rt.LoadTemplates()

	st := status.Template{T: t}
	st.LoadTemplates()

	router.SetHTMLTemplate(t)
}

func newOrg() *models.Organization {

	o := models.NewOrganization()
	o.OrgID = "org-oror"
	o.Alias = "oror"
	o.Emails = []models.Email{{Email: "email@email.email", Private: false}}
	o.Name = "Or Or Organization"
	o.OrgDetails.LocationStr = "Hyderabad"
	o.OrgDetails.LocationLL.Latitude = "17.235650"
	o.OrgDetails.LocationLL.Longitude = "79.124817"
	o.OrgDetails.Description = "string"
	return o
}
