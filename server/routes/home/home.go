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
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	// log.SetFlags(log.Llongfile | log.Ltime)
	log.Println("HERE")
	pkger.Include("/server/templates/")
	log.Println("HERE 2")
	home := router.Group("/home")
	{
		// router.LoadHTMLGlob("templates/*")
		tmpl, err := loadTemplate()
		log.Println("After load")
		if err != nil {
			log.Fatalln(err)
		}
		router.SetHTMLTemplate(tmpl)
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

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	var gblErr error
	log.Println("Load template")
	dd, err := pkger.Stat("/server/templates/")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(dd)
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
