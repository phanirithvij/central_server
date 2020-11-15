// +build !windows

package serve

import (
	"strconv"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func serve(router *gin.Engine, port int) {
	endless.ListenAndServe(":"+strconv.Itoa(port), router)
}
