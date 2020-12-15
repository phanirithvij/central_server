package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/phanirithvij/central_server/server/models"
	"github.com/phanirithvij/central_server/server/utils"
)

type orgSubmission struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenAuth token generation for orgs to get their own info
func TokenAuth(c *gin.Context) {
	sub := &orgSubmission{}
	err := c.BindJSON(sub)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	data := new(models.OrgSubmissionPass)
	data.Password = sub.Password
	data.OrgSubmission.Alias = sub.Email
	// go new doesn't create arrays
	data.OrgSubmission.Emails = []models.EmailD{{Email: sub.Email}}

	o, err := data.FindByEmail()
	if err != nil {
		// TODO check if banned or suspended or deleted
		// and send a custom message
		log.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":    err.Error(),
			"type":     "login",
			"messages": []string{"Couldn't find org with email " + sub.Email},
		})
		return
	}
	if utils.ComparePasswords(o.PasswordHash, data.Password) {
		// Save org id to cookie
		session := sessions.DefaultMany(c, "org")
		session.Set("org-id", o.ID)
		err = session.Save()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":    err.Error(),
				"type":     "cookie",
				"messages": []string{"Setting cookie failed"},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"type":     "login",
			"status":   "success",
			"messages": []string{"Authorized"},
		})
		return
	}
	c.JSON(http.StatusForbidden, gin.H{
		"error":    "Wrong password",
		"type":     "login",
		"messages": []string{"Your password is incorrect"},
	})
}
