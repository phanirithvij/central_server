// Package api the api that's exposed to the client apps and 3rd party client apps
package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/config"
	"github.com/phanirithvij/central_server/server/routes"
	v1 "github.com/phanirithvij/central_server/server/routes/api/v1"
	v2 "github.com/phanirithvij/central_server/server/routes/api/v2"
)

// SetupEndpoints Registers all the /api endpoints
func SetupEndpoints(router *gin.Engine) *gin.RouterGroup {
	apiG := router.Group("/api")
	{
		var versions []string
		v1gp := apiG.Group("/v1")
		{
			v1gp.GET("/", v1.Read)
			v1gp.GET("/read", v1.Read)

			orgrp := v1gp.Group("/orgs")
			{
				orgrp.GET("/:orgid/info", v1.OrgInfo)
			}

			home := v1gp.Group("/home")
			{
				home.GET("/public", v1.PublicList)
			}

			versions = append(versions, "v1")
		}
		v2gp := apiG.Group("/v2")
		{
			v2gp.GET("/", v2.Read)
			v2gp.GET("/read", v2.Read)
			// uncomment the next line when v2 api is ready
			// versions = append(versions, "v2")
		}
		apiG.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, apiG.BasePath()+"/versions")
		})
		apiG.GET("/versions", func(c *gin.Context) {
			c.String(http.StatusOK, strings.Join(versions, "\n"))
		})
	}
	routes.RegisterSelf(config.API)
	return apiG
}

func versionsRoute(versions []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, strings.Join(versions, "\n"))
	}
}
