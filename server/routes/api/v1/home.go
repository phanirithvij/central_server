package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phanirithvij/central_server/server/models"
	"gorm.io/gorm"
)

type urlParams struct {
	Pretty string `form:"pretty"`
}

// PublicList sends the public list of organizations
func PublicList(c *gin.Context) {
	org := models.NewOrganization()
	list, err := org.PublicList()
	msgs := []string{}
	status := http.StatusOK
	if err != nil {
		if errors.Is(err, models.ErrNoResultsFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			msgs = append(msgs, "No results found")
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{
			"error":    err.Error(),
			"status":   "failed",
			"messages": msgs,
		})
		return
	}
	p := &urlParams{}
	// https://github.com/gin-gonic/gin/issues/742
	c.BindQuery(p)
	if p.Pretty != "" {
		c.IndentedJSON(status, list)
		return
	}
	c.JSON(status, list)
}
