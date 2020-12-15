package v1

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type pingMessage struct {
	Message message `json:"message"`
}

type pongMessage struct {
	// just sending an ack
	Message string `json:"message"`
}

// Generally just `online`
// can also be `offline` if a planned maintanance or something
// in which case the server is marked offline with a special log entry
// indicating the server requested to be offline
type message string

// Ping the org servers ping the central server to let us know if they are online
// only then they'll be shown on the public lists
//
// The org servers should ping every 1 minute if they fail to do so in 2 minutes
// they will be marked offline and everyting's logged into their activity dashboard
// they can ping whenever after that and they'll be made online
// TODO if caching is implemented it should take this online status thing into account
func Ping(c *gin.Context) {
	ping := &pingMessage{}
	err := c.BindJSON(ping)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "JSON parsing failed",
		})
		return
	}
	session := sessions.DefaultMany(c, "org")
	v, ok := session.Get("org-id").(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":    "session has no org ID",
			"type":     "login",
			"messages": []string{"Not Authorized"},
		})
		return
	}
	if ping.Message == "ping" {
		// TODO set status online
		// log to activity
		log.Println("Org server for", v, "is online")
	} else if ping.Message == "offline" {
		// TODO set status offline
		// log to activity as an offline request
		log.Println("Org server for", v, "is going offline")
	}
	log.Println(ping.Message)
	pong := &pongMessage{"pong"}
	c.JSON(http.StatusOK, pong)
}
