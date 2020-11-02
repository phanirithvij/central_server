package home

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	"github.com/phanirithvij/central_server/server/utils"
)

//go:generate pkger -o server

// Called checks whether we called it or not
var Called = false

// PkgerPrefix the prefix and the top level dir for all the assets
const (
	PkgerPrefix      = "github.com/phanirithvij/central_server:"
	HomeAssetsPrefix = PkgerPrefix + `/server/routes/home/`
)

func init() {
	pkger.Include("/server/routes/home/home.html")
}

// RegisterEndPoints Registers all the /api endpoints
// Must call LoadTemplates before this if it exists
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	// log.SetFlags(log.Llongfile | log.Ltime)
	log.Println("Called", Called)
	if !Called {
		log.Println("[Warning] templates are uninitialized for home route")
	}
	home := router.Group("/home")
	{
		home.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, HomeAssetsPrefix+"home.html", gin.H{
				"title": "Home Page",
			})
		})
		home.GET("/hello", func(c *gin.Context) {
			c.String(http.StatusOK, `strings.Join(versions, "\n")`)
		})
	}
	return home
}

// https://gin-gonic.com/docs/examples/bind-single-binary-with-template/
// https://github.com/gin-gonic/examples/commit/c5a87f03d39fdb9e0f6312344c21ccdd55140293

// LoadTemplates loads the templates used by this package
func LoadTemplates(t *template.Template) (*template.Template, error) {
	Called = true
	return utils.LoadTemplates(t, HomeAssetsPrefix)
}
