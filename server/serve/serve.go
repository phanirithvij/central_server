// Package serve serves the server
package serve

import (
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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
	assetDir  = "/client/vue/build"
)

func init() {
	pkger.Include("/client/vue/build")
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

	// https://github.com/gorilla/mux#serving-single-page-applications
	spa := &spaHandler{
		staticPath: assetDir,
		indexPath:  assetDir + "/index.html",
	}

	// https://stackoverflow.com/a/34373030/8608146
	gzHandler := gziphandler.GzipHandler(spa)
	cacheH := http.StripPrefix(fbBaseURL, cache(gzHandler))
	router.GET(fbBaseURL+"/*w", gin.WrapH(cacheH))

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

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	// TODO check if path has .. i.e relative routes and ban IP
	// https://github.com/mrichman/godnsbl
	// https://github.com/jpillora/ipfilter

	// prepend the path with the path to the static directory
	path := filepath.Join(h.staticPath, r.URL.Path)

	// check whether a file exists at the given path
	_, err := pkger.Stat(path)
	if err != nil {
		// file does not exist, serve index.html
		// http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		file, err := pkger.Open(h.indexPath)
		if err != nil {
			http.Error(w, "file "+r.URL.Path+" does not exist", http.StatusNotFound)
			return
		}
		// lw := lhWriter{w}
		lw := w
		// r.URL.Path += "/index.html"
		cont, err := ioutil.ReadAll(file)
		lw.Header().Set("Content-Type", "text/html")
		lw.Write(cont)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(pkger.Dir(h.staticPath)).ServeHTTP(w, r)
}

var (
	// server start time
	serverStart    = time.Now()
	serverStartStr = serverStart.Format(http.TimeFormat)
	expireDur      = time.Minute * 10
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
		fi, err := pkger.Stat(filepath.Join(assetDir, fname))

		if err != nil {
			// spa route eg. /web/about
			// let spa handle it, no need to cache
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

type lhWriter struct {
	w http.ResponseWriter
}

func (w lhWriter) Write(b []byte) (int, error) {
	log.Println(string(b))
	return w.w.Write(b)
}

func (w lhWriter) WriteHeader(code int) {
	w.w.WriteHeader(code)
}

func (w lhWriter) Header() http.Header {
	return w.w.Header()
}
