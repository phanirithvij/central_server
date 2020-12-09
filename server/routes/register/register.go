// Package register Home routes
package register

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
	dbm "github.com/phanirithvij/central_server/server/utils/db"
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

// TemplateParams for this route
type TemplateParams struct {
	Title string
}

// SetupEndpoints Registers all the /register endpoints
// Must call LoadTemplates before this if it exists
// Returns the router group so it can be also used to set routes externally
func SetupEndpoints(router *gin.Engine) *gin.RouterGroup {
	db := dbm.GetDB()
	if !templatesInitDone {
		log.Fatalln(errors.New(usage))
	}
	register := router.Group("/register")
	{
		register.GET("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			// TODO get the currently loggedin orgid
			// then get it from db
			session := sessions.DefaultMany(c, "org")
			_, ok := session.Get("org-id").(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":    "session has no org ID",
					"type":     "register",
					"messages": []string{"Not Authorized"},
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"type":     "register",
				"status":   "success",
				"messages": []string{"Authorized"},
			})
		})
		register.OPTIONS("/*_", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Methods", "POST, GET")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			c.Header("Access-Control-Allow-Credentials", "true")

			c.Status(http.StatusOK)
		})
		register.POST("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			d := json.NewDecoder(c.Request.Body)
			data := &models.OrgSubmissionPass{}
			err := d.Decode(&data)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			// o := data.
			o := models.NewOrganization()
			o.Alias = data.Alias
			// TODO all emails
			// TODO frontend email private option
			// TODO don't ask all these details when signing up
			// Ask after email verification
			// TODO multi email verification
			o.Emails = []models.Email{}
			for _, email := range data.Emails {
				o.Emails = append(o.Emails, models.Email{
					Email:   email.Email,
					Private: email.Private,
				})
			}
			o.Name = data.Name
			msgs, err := o.ValidateSub([]string{"Name", "Emails", "Alias"})
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
			session := sessions.DefaultMany(c, "org")
			session.Set("org-id", o.ID)
			err = session.Save()
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":    err.Error(),
					"type":     "cookie",
					"messages": []string{"Setting cookie failed"},
				})
				return
			}
			c.JSON(http.StatusCreated, o)
		})
		// check if alias exists in database
		register.GET("/alias/:alias", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			type alias struct {
				Alias string
			}
			al := c.Param("alias")
			o := alias{}
			if err := db.
				Model(&models.Organization{}).
				Where("alias = ?", al).
				Find(&o).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": err.Error()})
				return
			}
			if o.Alias != "" {
				c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Alias " + al + " already exists"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Alias " + al + " avaliable"})
		})
	}
	routes.RegisterSelf(config.Register)
	logout := router.Group("/logout")
	{
		logout.GET("/", func(c *gin.Context) {
			// logout will remove session
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			session := sessions.DefaultMany(c, "org")
			_, ok := session.Get("org-id").(uint)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":    "session has no org ID",
					"type":     "logout",
					"messages": []string{"Not Authorized"},
				})
				return
			}
			session.Set("org-id", nil)
			err := session.Save()
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":    err.Error(),
					"type":     "logout",
					"messages": []string{"Couldn't clear session"},
				})
			}
			c.JSON(http.StatusAccepted, gin.H{
				"type":     "logout",
				"messages": []string{"Logged out"},
			})

		})
	}
	routes.RegisterSelf(config.Logout)
	return register
}

// Template a wrapper of template.Template
type Template struct {
	T *template.Template
}

// LoadTemplates loads the templates used by register package
func (t Template) LoadTemplates() {
	templatesInitDone = true
}
