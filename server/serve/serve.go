// Package serve serves the server
package serve

import (
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NYTimes/gziphandler"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/phanirithvij/central_server/server/models"
	"github.com/phanirithvij/central_server/server/routes"
	api "github.com/phanirithvij/central_server/server/routes/api"
	home "github.com/phanirithvij/central_server/server/routes/home"
	register "github.com/phanirithvij/central_server/server/routes/register"
	status "github.com/phanirithvij/central_server/server/routes/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	fbBaseURL = "/web"
	assetDir  = "/client/web/build"
)

func init() {
	pkger.Include("/client/web/build")
}

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
	// https://github.com/gin-gonic/gin/issues/293#issuecomment-103659145
	gzHandler := gziphandler.GzipHandler(http.FileServer(pkger.Dir(assetDir)))
	statik := http.StripPrefix(fbBaseURL, cache(gzHandler))
	router.GET(fbBaseURL+"/*w", gin.WrapH(statik))

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
	serverStart    = time.Now()
	serverStartStr = serverStart.Format(http.TimeFormat)
	expireDur      = time.Minute * 2
	expire         = serverStart.Add(expireDur)
	expireStr      = expire.Format(http.TimeFormat)
)

// https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702

// cache caching the public directory
func cache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fname := r.URL.Path
		if r.URL.Path == fbBaseURL {
			fname = "index.html"
		}
		fi, err := pkger.Stat(assetDir + fname)

		if err != nil {
			log.Println(err)
			h.ServeHTTP(w, r)
			return
		}
		modTime := fi.ModTime()

		fhex := hex.EncodeToString([]byte(fname))
		fmodTH := hex.EncodeToString([]byte(strconv.FormatInt(modTime.Unix(), 10)))
		etagH := fhex + "." + fmodTH

		etag := r.Header.Get("If-None-Match")
		if etag != "" && etag == etagH {
			w.WriteHeader(http.StatusNotModified)
			w.Header().Set("Cache-Control", "public, max-age="+strconv.FormatInt(int64(expireDur.Seconds()), 10))
			w.Header().Set("Expires", expireStr)
			w.Header().Set("Etag", etagH)
			return
		}

		// https://stackoverflow.com/a/48876760/8608146

		w.Header().Set("Cache-Control", "public, max-age="+strconv.FormatInt(int64(expireDur.Seconds()), 10))
		w.Header().Set("Expires", expireStr)
		w.Header().Set("Etag", etagH)

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
