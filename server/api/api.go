package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	v1 "github.com/phanirithvij/central_server/server/api/v1"
	v2 "github.com/phanirithvij/central_server/server/api/v2"
)

// RegisterEndPoints Registers all the /api endpoints
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	apiG := router.Group("/api")
	{
		var versions []string
		v1gp := apiG.Group("/v1")
		{
			v1gp.GET("/read", v1.Read)
			// v1gp.POST("/read", v1.Read)
			// v1gp.DELETE("/read", v1.Read)
			// v1gp.PATCH("/read", v1.Read)
			// v1gp.PUT("/read", v1.Read)
			versions = append(versions, "v1")
		}
		v2gp := apiG.Group("/v2")
		{
			v2gp.GET("/read", v2.Read)
			// v2gp.POST("/read", v2.Read)
			// v2gp.DELETE("/read", v2.Read)
			// v2gp.PATCH("/read", v2.Read)
			// v2gp.PUT("/read", v2.Read)
			versions = append(versions, "v2")
		}
		apiG.GET("/versions", func(c *gin.Context) {
			c.String(http.StatusOK, strings.Join(versions, "\n"))
		})
	}
	return apiG
}
