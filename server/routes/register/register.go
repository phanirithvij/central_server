// Package register Home routes
package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/config"
	"github.com/phanirithvij/central_server/server/models"
	"github.com/phanirithvij/central_server/server/routes"
	"github.com/phanirithvij/central_server/server/utils"
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
	register := router.Group("/apiOrg/register")
	{
		optsCors := cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
			AllowMethods:     []string{"POST"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
			AllowCredentials: true,
		})
		register.OPTIONS("/*_", optsCors)

		credCors := cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
			AllowCredentials: true,
		})

		// Enable CORS for react client when in dev
		register.POST("/", credCors, func(c *gin.Context) {
			d := json.NewDecoder(c.Request.Body)
			data := &models.OrgSubmissionPass{}
			err := d.Decode(&data)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			// TODO frontend email private option
			// TODO don't ask all these details when signing up
			// Ask after email verification
			// TODO multi email verification
			o := data.Org()
			// set this is as the mail email
			True := true
			o.Emails[0].Main = &True
			msgs, err := o.ValidateSub([]string{"Emails", "Alias"})
			if err != nil {
				log.Println(err)
				// https://stackoverflow.com/a/40926661/8608146
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error":    err.Error(),
					"type":     "validate",
					"messages": msgs,
				})
				return
			}
			o.PasswordHash = utils.Hash(data.Password)
			log.Println(data)
			err = o.SaveReq(c)
			if err != nil {
				log.Println(err)
				return
			}
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
			c.JSON(http.StatusCreated, o)
		})
		// check if alias exists in database
		register.GET("/alias/:alias", credCors, func(c *gin.Context) {
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
