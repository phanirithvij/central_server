package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/models"
)

type urlParams struct {
	Pretty string `form:"pretty"`
}

// PublicList sends the public list of organizations
func PublicList(c *gin.Context) {
	org := models.NewOrganization()
	list, err := org.PublicList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err.Error(),
			"status":   "failed",
			"messages": []string{"PublicList could not be retreived"},
		})
		return
	}
	p := &urlParams{}
	// https://github.com/gin-gonic/gin/issues/742
	c.BindQuery(p)
	if p.Pretty != "" {
		c.IndentedJSON(http.StatusOK, list)
		return
	}
	c.JSON(http.StatusOK, list)
}
