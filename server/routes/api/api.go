// Package api the api that's exposed to the client apps and 3rd party client apps
package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	v1 "github.com/phanirithvij/central_server/server/routes/api/v1"
	v2 "github.com/phanirithvij/central_server/server/routes/api/v2"
)

var (
	// EndpointsRegistered to keep track of whether the endpoints are registered for this route
	EndpointsRegistered = false
)

// RegisterEndPoints Registers all the /api endpoints
func RegisterEndPoints(router *gin.Engine) *gin.RouterGroup {
	apiG := router.Group("/api")
	{
		var versions []string
		v1gp := apiG.Group("/v1")
		{
			v1gp.GET("/", v1.Read)
			v1gp.GET("/read", v1.Read)
			versions = append(versions, "v1")
		}
		v2gp := apiG.Group("/v2")
		{
			v2gp.GET("/", v2.Read)
			v2gp.GET("/read", v2.Read)
			versions = append(versions, "v2")
		}
		apiG.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, apiG.BasePath()+"/versions")
		})
		apiG.GET("/versions", func(c *gin.Context) {
			c.String(http.StatusOK, strings.Join(versions, "\n"))
		})
	}
	EndpointsRegistered = true
	return apiG
}

func versionsRoute(versions []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, strings.Join(versions, "\n"))
	}
}
