package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
	userservice "github.com/kirtfieldk/astella/src/services/userService"
)

func (h *BaseHandler) GeteventsMemberOf(c *gin.Context) {
	userId := c.Param(constants.USER_ID)
	events, err := userservice.GetUserEvents(userId, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, events)

}
