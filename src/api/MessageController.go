package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
	messageservice "github.com/kirtfieldk/astella/src/services/messageService"
	"github.com/kirtfieldk/astella/src/structures"
)

func (h *BaseHandler) PostMessageToEvent(c *gin.Context) {
	var msg structures.Message
	if err := c.BindJSON(&msg); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.PAYLOAD_IS_NOT_MSG})
		return
	}

	success, err := messageservice.PostMessage(msg, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{constants.MESSAGE: success})
}
