// Package serve serves the server
package serve

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/phanirithvij/central_server/server/models"
	"github.com/phanirithvij/central_server/server/routes"
	api "github.com/phanirithvij/central_server/server/routes/api"
	home "github.com/phanirithvij/central_server/server/routes/home"
	register "github.com/phanirithvij/central_server/server/routes/register"
	status "github.com/phanirithvij/central_server/server/routes/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const fbBaseURL = "/web"

//go:generate packr2

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

	// https://stackoverflow.com/a/55854101/8608146
	// router.Static("/web", "./client/web/build")
	box := packr.New(fbBaseURL, "../../client/web/build")
	// router.StaticFS("/web", box)
	// https://github.com/gin-gonic/gin/issues/293#issuecomment-103659145
	router.Any(fbBaseURL+"*", gin.WrapH(cache(http.FileServer(box))))

	o := newOrg()
	// utils.PrintStruct(*o)
	// o.Print()
	o.Validate()

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Organization{}, &models.Email{})
	if err != nil {
		log.Fatalln(err)
	}

	db.Create(&o)
	db.Save(&o)

	serve(router, port)
}

var (
	// server start time
	cacheSince    = time.Now()
	cacheSinceStr = cacheSince.Format(http.TimeFormat)
	cacheUntil    = cacheSince.AddDate(0, 0, 7)
	cacheUntilStr = cacheUntil.Format(http.TimeFormat)
)

// https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702

// cache caching the public directory
func cache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, fbBaseURL+"/static") {
			modtime := r.Header.Get("If-Modified-Since")
			if modtime != "" {
				// TODO check if file is modified since time `t`
				// t, _ := time.Parse(http.TimeFormat, modtime)
				log.Println("[Warning] not checking if modified")
				w.WriteHeader(http.StatusNotModified)
				// no need to forward as cache
				return
			}

			// 604800 -> one week
			w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
			w.Header().Set("Last-Modified", cacheSinceStr)
			w.Header().Set("Expires", cacheUntilStr)
		}
		// forward
		h.ServeHTTP(w, r)
	})
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
