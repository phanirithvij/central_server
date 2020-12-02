package serve

import (
	"log"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

// On window endless library will not build so we use normal http server
func serve(router *gin.Engine, port int) {
	log.Println("Serving on Port", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}
