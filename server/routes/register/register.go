// Package register Home routes
package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/phanirithvij/central_server/server/config"
	"github.com/phanirithvij/central_server/server/models"
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
			o := models.NewOrganization()
			o.Alias = data.Alias
			// TODO all emails
			// TODO frontend email private option
			// TODO don't ask all these details when signing up
			// Ask after email verification
			// TODO multi email verification
			o.Emails = []models.Email{{Email: data.Emails[0], Private: false}}
			o.Name = data.Name
			o.OrgDetails.LocationStr = data.Address
			o.OrgDetails.LocationLL.Latitude = strconv.FormatFloat(data.Location[0], 'f', -1, 64)
			o.OrgDetails.LocationLL.Longitude = strconv.FormatFloat(data.Location[1], 'f', -1, 64)
			o.OrgDetails.Description = data.Description
			msgs, err := o.Validate()
			if err != nil {
				// https://stackoverflow.com/a/40926661/8608146
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":    err.Error(),
					"type":     "validate",
					"messages": msgs,
				})
				return
			}
			log.Println(data)
			log.Println(o.Str())
			err = o.SaveReq(c)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(o.Str())
			c.JSON(http.StatusOK, o)
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
