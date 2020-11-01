package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// VersionInfo API version Info for v1
type VersionInfo struct {
	Version string `json:"version"`
}

// Read the read endpoint for api.v1
func Read(c *gin.Context) {
	resp := VersionInfo{Version: "v1"}
	c.JSON(http.StatusOK, gin.H{
		"info": resp,
	})
}
