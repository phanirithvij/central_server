package serve

import (
	"html/template"
	"log"
	"strconv"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	api "github.com/phanirithvij/central_server/server/api"
	home "github.com/phanirithvij/central_server/server/routes/home"
)

// Serve A function which serves the server
func Serve(port int, debug bool) {
	if debug {
		log.SetFlags(log.Ltime | log.Llongfile)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	t := template.New("")

	t, err := home.LoadTemplates(t)
	if err != nil {
		log.Fatalln(err)
	}

	router.SetHTMLTemplate(t)

	api.RegisterEndPoints(router)
	home.RegisterEndPoints(router)

	endless.ListenAndServe(":"+strconv.Itoa(port), router)
}
