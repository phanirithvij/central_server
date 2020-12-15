// Package v1 Api version v1
package v1

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/models"
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

// OrgInfo sends the organization info
func OrgInfo(c *gin.Context) {
	// TODO: Allow cors for the server currently sending the request
	// if it's registered in the DB
	org := models.NewOrganization()

	session := sessions.DefaultMany(c, "org")
	v, ok := session.Get("org-id").(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":    "session has no org ID",
			"type":     "settings",
			"messages": []string{"Not Authorized"},
		})
		return
	}

	// https://stackoverflow.com/questions/16427416/what-are-the-advantages-of-the-general-types-int-uint-over-specific-types-i#comment23559803_16427485
	org.ID = v
	org, err := org.Find()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":    err.Error(),
			"status":   "failed",
			"messages": []string{"Organization not found"},
		})
		return
	}
	c.JSON(http.StatusOK, org.OrgSubmission())
}
