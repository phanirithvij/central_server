package home

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
)

//go:generate pkger -o server

// RegisterEndPoints Registers all the /api endpoints
// Must call LoadTemplates before this if it exists
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	// log.SetFlags(log.Llongfile | log.Ltime)
	log.Println("HERE")
	pkger.Include("/server/templates/")
	log.Println("HERE 2")
	home := router.Group("/home")
	{
		// router.LoadHTMLGlob("templates/*")
		// router.LoadHTMLFiles("home.html")
		home.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"title": "Home Page",
			})
		})
		home.GET("/hello", func(c *gin.Context) {
			c.String(http.StatusOK, `strings.Join(versions, "\n")`)
		})
	}
	return home
}

// LoadTemplates loads the templates used by this package
func LoadTemplates(t *template.Template) (*template.Template, error) {
	var gblErr error
	gblErr = pkger.Walk("/server/templates/", func(path string, info os.FileInfo, err error) error {
		log.Println(path, err)
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		file, err := pkger.Open(path)
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		t, err = t.New(path).Parse(string(h))
		if err != nil {
			return err
		}
		return nil
	})
	if gblErr != nil {
		return nil, gblErr
	}
	return t, nil
}
