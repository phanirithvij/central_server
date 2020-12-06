// Package register Home routes
package register

import (
	"encoding/json"
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
)

var (
	// to keep track of whether the templates are initialized or not for this route
	templatesInitDone = false
)

// PkgerPrefix the prefix and the top level dir for all the assets
const (
	registerAssetsPrefix = config.PkgerPrefix + `/server/routes/` + config.Register
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

`, config.Register)

func init() {
	// include dirs for pkger parser to pickup
	pkger.Include("/server/routes/register/register.html")
}

// TemplateParams for this route
type TemplateParams struct {
	Title string
}

type orgSubmission struct {
	Address     string    `json:"address"`
	Alias       string    `json:"alias"`
	Description string    `json:"description"`
	Emails      []string  `json:"emails"`
	Location    []float64 `json:"location"`
	Name        string    `json:"name"`
	Password    string    `json:"password"`
}

// RegisterEndPoints Registers all the /register endpoints
// Must call LoadTemplates before this if it exists
// Returns the router group so it can be also used to set routes externally
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	if !templatesInitDone {
		log.Fatalln(errors.New(usage))
	}
	register := router.Group("/register")
	{
		register.GET("/", func(c *gin.Context) {
			params := TemplateParams{Title: "register page"}
			c.HTML(http.StatusOK, registerAssetsPrefix+"/register.html", params)
		})
		register.OPTIONS("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Methods", "POST")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

			c.Status(http.StatusOK)
		})
		register.POST("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			d := json.NewDecoder(c.Request.Body)
			data := &orgSubmission{}
			err := d.Decode(&data)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			log.Println(data)
			c.JSON(http.StatusOK, data)
		})
		register.GET("/hello", func(c *gin.Context) {
			c.String(http.StatusOK, `strings.Join(versions, "\n")`)
		})
	}
	routes.RegisterSelf(config.Register)
	return register
}

// Template a wrapper of template.Template
type Template struct {
	T *template.Template
}

// LoadTemplates loads the templates used by register package
func (t Template) LoadTemplates() {
	before := len(t.T.Templates())
	_, err := utils.LoadTemplates(t.T, registerAssetsPrefix)
	if err != nil {
		log.Fatalln(err)
	}
	after := len(t.T.Templates())
	if before < after {
		templatesInitDone = true
	}
}
