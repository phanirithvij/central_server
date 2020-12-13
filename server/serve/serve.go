// Package serve serves the server
package serve

import (
	"compress/gzip"
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/NYTimes/gziphandler"

	"github.com/didip/tollbooth/v6"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/phanirithvij/central_server/server/models"
	"github.com/phanirithvij/central_server/server/routes"
	api "github.com/phanirithvij/central_server/server/routes/api"
	home "github.com/phanirithvij/central_server/server/routes/home"
	login "github.com/phanirithvij/central_server/server/routes/login"
	register "github.com/phanirithvij/central_server/server/routes/register"
	settings "github.com/phanirithvij/central_server/server/routes/settings"
	status "github.com/phanirithvij/central_server/server/routes/status"
	"github.com/phanirithvij/central_server/server/utils"
	dbm "github.com/phanirithvij/central_server/server/utils/db"
	"github.com/phanirithvij/central_server/server/utils/rate"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	orgBaseURL    = "/org"
	adminBaseURL  = "/admin"
	orgAssetDir   = "/client/org/build"
	adminAssetDir = "/client/admin/build"
)

func init() {
	pkger.Include("/client/org/build")
	pkger.Include("/client/admin/build")
}

// Serve A function which serves the server
func Serve(port int, debug bool) {
	if debug {
		log.SetFlags(log.Ltime | log.Lshortfile)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	dbm.InitDB(debug)
	db := dbm.DB

	o := newOrg()
	// utils.PrintStruct(*o)
	// o.Print()
	o.Validate()

	// Migrate the schema
	err := db.AutoMigrate(&models.Organization{}, &models.Email{}, &models.Server{})
	if err != nil {
		log.Fatalln(err)
	}

	db.Create(o)
	// db.Save(o)

	router := gin.Default()
	registerTemplates(router)
	setupSessionStore(router)

	api.SetupEndpoints(router)
	home.SetupEndpoints(router)
	register.SetupEndpoints(router)
	login.SetupEndpoints(router)
	status.SetupEndpoints(router)
	settings.SetupEndpoints(router)

	routes.CheckEndpoints()

	// https://stackoverflow.com/a/55854101/8608146
	// https://github.com/gin-gonic/gin/issues/293#issuecomment-103659145
	// https://create-react-app.dev/docs/adding-custom-environment-variables/

	// https://github.com/gorilla/mux#serving-single-page-applications
	orgSPA := &spaHandler{
		staticPath: orgAssetDir,
		indexPath:  orgAssetDir + "/index.html",
	}

	gh, err := gziphandler.NewGzipLevelHandler(gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	orgGzHandler := gh(orgSPA)
	orgCacheH := http.StripPrefix(orgBaseURL, cache(orgGzHandler, orgAssetDir))
	router.GET(orgBaseURL+"/*w", gin.WrapH(orgCacheH))

	amdinSPA := &spaHandler{
		staticPath: adminAssetDir,
		indexPath:  adminAssetDir + "/index.html",
	}

	adminGzHandler := gh(amdinSPA)
	adminCacheH := http.StripPrefix(adminBaseURL, cache(adminGzHandler, adminAssetDir))
	router.GET(adminBaseURL+"/*w", gin.WrapH(adminCacheH))

	promH := promhttp.Handler()
	lmt := tollbooth.NewLimiter(3, nil)
	rateF := rate.LimitHandler(lmt)
	router.GET("/metrics", rateF, gin.WrapH(promH))

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
	// github.com/didip/tollbooth/v6

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
func cache(h http.Handler, assetDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fname := r.URL.Path
		if r.URL.Path == orgBaseURL {
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

	lt := login.Template{T: t}
	lt.LoadTemplates()

	st := status.Template{T: t}
	st.LoadTemplates()

	set := settings.Template{T: t}
	set.LoadTemplates()

	router.SetHTMLTemplate(t)
}

// True true pointer
func True() *bool {
	True := true
	return &True
}

// False False pointer
func False() *bool {
	False := false
	return &False
}

func newOrg() *models.Organization {
	o := models.NewOrganization()
	o.PasswordHash = utils.Hash("oror")
	// o.OrgID = "org-oror"

	o.Private = False()
	o.Alias = "oror"
	o.Emails = []models.Email{
		{Email: "emaixl@email.emailemail", Private: False(), Main: True()},
		{Email: "email3w@email3.email", Private: True()},
		{Email: "emai2lw@x.email", Private: False()},
		{Email: "emxaxi2lwx@email.", Private: True()},
		{Email: "emxaxilwx@xxemail.email", Private: False()},
		{Email: "emailwxxw@exmxail.email", Private: True()},
		{Email: "emailw@wemaixl.email", Private: False()},
	}

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

func setupSessionStore(r *gin.Engine) sessions.Store {
	secretKey := os.Getenv("SESSION_SECRET")
	if secretKey == "" {
		// https://randomkeygen.com/
		secretKey = "$l9voQ>MoLq{nAT#zJzH*b_;jC=2g6"
	}
	bsk := []byte(secretKey)
	// https://github.com/gin-contrib/sessions#redis
	// https://github.com/boj/redistore/blob/cd5dcc76aeff9ba06b0a924829fe24fd69cdd517/redistore.go#L155
	// size: maximum number of idle connections.
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", bsk)
	if err != nil {
		log.Println(err)
		log.Println("[WARNING] Redis not available so using cookie sessions")
		store = cookie.NewStore(bsk)
	}
	// https://github.com/gin-contrib/sessions#multiple-sessions
	sessionNames := []string{"org", "admin"}
	r.Use(sessions.SessionsMany(sessionNames, store))
	return store
}
