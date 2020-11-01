package serve

import (
	"strconv"

	"github.com/gin-gonic/gin"
	api "github.com/phanirithvij/central_server/server/api"
	home "github.com/phanirithvij/central_server/server/routes/home"
)

// Serve A function which serves the server
func Serve(port int) {
	router := gin.Default()

	api.RegisterEndPoints(router)
	home.RegisterEndPoints(router)

	router.Run(":" + strconv.Itoa(port))
}
