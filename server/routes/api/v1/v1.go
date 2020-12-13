// Package v1 Api version v1
package v1

import (
	"net/http"
	"strconv"

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
	orgid, ok := c.Params.Get("orgid")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error":  "Org id was not found in the url",
			"status": "failed",
		})
		return
	}
	orgIduint, err := strconv.ParseUint(orgid, 10, 64)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error":  "Org id was not a valid uint in the url",
			"status": "failed",
		})
		return
	}

	// https://stackoverflow.com/questions/16427416/what-are-the-advantages-of-the-general-types-int-uint-over-specific-types-i#comment23559803_16427485
	org.ID = uint(orgIduint)
	org, err = org.Find()
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
