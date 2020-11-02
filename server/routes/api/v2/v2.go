// Package v2 Api version 2
package v2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// VersionInfo API version Info for v2
type VersionInfo struct {
	Version string `json:"version"`
}

// Read the read endpoint for api.v2
func Read(c *gin.Context) {
	resp := VersionInfo{Version: "v2"}
	// breaking response format change
	c.JSON(http.StatusOK, resp)
}
