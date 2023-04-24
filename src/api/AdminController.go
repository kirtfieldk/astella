package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Admin struct {
	UserId  string `json:"user_id"`
	EventId string `json:"event_id"`
}

func addAdmin(c *gin.Context) {
	var newAdmin Admin
	if err := c.BindJSON(&newAdmin); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, newAdmin)

}
