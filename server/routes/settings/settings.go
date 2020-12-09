// Package settings Home routes
package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/config"
	"github.com/phanirithvij/central_server/server/models"
	"github.com/phanirithvij/central_server/server/routes"
)

var (
	// to keep track of whether the templates are initialized or not for this route
	templatesInitDone = false
)

// PkgerPrefix the prefix and the top level dir for all the assets
const (
	settingsAssetsPrefix = config.PkgerPrefix + `/server/routes/` + config.Settings
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

`, config.Settings)

// TemplateParams for this route
type TemplateParams struct {
	Title string
}

// SetupEndpoints settingss all the /settings endpoints
// Must call LoadTemplates before this if it exists
// Returns the router group so it can be also used to set routes externally
func SetupEndpoints(router *gin.Engine) *gin.RouterGroup {
	if !templatesInitDone {
		log.Fatalln(errors.New(usage))
	}
	settings := router.Group("/settings")
	{
		settings.GET("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			// TODO get the currently loggedin orgid
			// then get it from db
			session := sessions.DefaultMany(c, "org")
			data := &models.OrgSubmission{}
			v, ok := session.Get("org-id").(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":    "session has no org ID",
					"type":     "settings",
					"messages": []string{"Not Authorized"},
				})
				return
			}
			data.ID = v
			o, err := data.Find()
			// deleted organizations with active sessions
			if o.Alias == "" {
				// Alias is empty => not found on server
				err = errors.New("Organization not found")
			}
			log.Println(o.Str())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error":    err.Error(),
					"type":     "settings",
					"messages": []string{"Organization not found"},
				})
				return
			}
			c.JSON(http.StatusOK, o.OrgSubmission())
		})

		settings.OPTIONS("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Methods", "PUT")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			c.Header("Access-Control-Allow-Credentials", "true")

			c.Status(http.StatusOK)
		})

		settings.PUT("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			// cookies won't show in react devtools
			// https://stackoverflow.com/a/50370345/8608146
			// log.Println(c.Request.Header.Get("Cookie"))
			d := json.NewDecoder(c.Request.Body)
			data := &models.OrgSubmission{}
			err := d.Decode(&data)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":    err.Error(),
					"type":     "json",
					"messages": []string{"Server got an invalid JSON"},
				})
				return
			}
			session := sessions.DefaultMany(c, "org")
			v, ok := session.Get("org-id").(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":    "sesson has no org ID",
					"type":     "settings",
					"messages": []string{"Not Authorized"},
				})
				return
			}
			data.ID = v
			o := data.Org()
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusNotFound, gin.H{
					"error":    err.Error(),
					"type":     "no-org",
					"messages": []string{"No such organization found"},
				})
				return
			}
			// Allow update only after email verification
			// TODO multi email verification
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
			err = o.SaveReq(c)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(o.Str())
			c.JSON(http.StatusOK, o)
		})
	}
	routes.RegisterSelf(config.Settings)
	return settings
}

// Template a wrapper of template.Template
type Template struct {
	T *template.Template
}

// LoadTemplates loads the templates used by settings package
func (t Template) LoadTemplates() {
	templatesInitDone = true
}
