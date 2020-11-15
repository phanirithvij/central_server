package serve

import (
	"log"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

func serve(router *gin.Engine, port int) {
	log.Println("Serving on Port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), router)
}
