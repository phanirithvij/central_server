// Package login Home routes
package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

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
	loginAssetsPrefix = config.PkgerPrefix + `/server/routes/` + config.Login
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

`, config.Login)

// TemplateParams for this route
type TemplateParams struct {
	Title string
}

type loginSubmission struct {
	EmailAlias string `json:"emailAlias"`
	Password   string `json:"password"`
}

// SetupEndpoints Registers all the /login endpoints
// Must call LoadTemplates before this if it exists
// Returns the router group so it can be also used to set routes externally
func SetupEndpoints(router *gin.Engine) *gin.RouterGroup {
	_ = dbm.GetDB()
	if !templatesInitDone {
		log.Fatalln(errors.New(usage))
	}
	login := router.Group("/login")
	{
		login.GET("/", func(c *gin.Context) {
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
					"type":     "login",
					"messages": []string{"Not Authorized"},
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"type":     "login",
				"status":   "success",
				"messages": []string{"Authorized"},
			})
		})
		login.OPTIONS("/*_", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Methods", "POST, GET")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			c.Header("Access-Control-Allow-Credentials", "true")

			c.Status(http.StatusOK)
		})
		login.POST("/", func(c *gin.Context) {
			// Enable CORS for react client when in dev
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			d := json.NewDecoder(c.Request.Body)
			sub := &loginSubmission{}
			err := d.Decode(&sub)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			data := new(models.OrgSubmissionPass)
			data.Password = sub.Password
			data.OrgSubmission.Alias = sub.EmailAlias
			// go new doesn't create arrays
			data.OrgSubmission.Emails = []models.EmailD{{Email: sub.EmailAlias}}

			const AliasX = 0
			const EmailX = 1
			const NoMethodX = 2
			method := EmailX

			o := data.Org()
			// TODO see if any one of these is valid and use that
			_, err = o.ValidateSub([]string{"Alias", "Emails"})
			log.Println(err)
			if err != nil {
				if strings.Contains(err.Error(), "Email") {
					// not a valid email try alias
					method = AliasX
					if strings.Contains(err.Error(), "Alias") {
						// invalid alias
						method = NoMethodX
					}
				}
			}

			log.Println(data)
			log.Println(method)
			log.Println(o.Str())

			switch method {
			case AliasX:
				o, err = data.FindByAlias()
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"error":    err.Error(),
						"type":     "login",
						"messages": []string{"Couldn't find org with alias " + sub.EmailAlias},
					})
					return
				}
			case EmailX:
				o, err = data.FindByEmail()
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"error":    err.Error(),
						"type":     "login",
						"messages": []string{"Couldn't find org with email " + sub.EmailAlias},
					})
					return
				}
			default:
				// validation failed for both email and alias
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":    err.Error(),
					"type":     "validate",
					"messages": []string{sub.EmailAlias + " is not a valid email address or an alias"},
				})
				return
			}

			log.Println(o.Str())

			// Save org id to cookie
			session := sessions.DefaultMany(c, "org")
			session.Set("org-id", o.ID)
			err = session.Save()
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":    err.Error(),
					"type":     "cookie",
					"messages": []string{"Setting cookie failed"},
				})
				return
			}
			c.JSON(http.StatusOK, o)
		})
	}
	routes.RegisterSelf(config.Login)
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
	return login
}

// Template a wrapper of template.Template
type Template struct {
	T *template.Template
}

// LoadTemplates loads the templates used by login package
func (t Template) LoadTemplates() {
	templatesInitDone = true
	// not using any templates
}