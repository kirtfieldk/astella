package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
	messageservice "github.com/kirtfieldk/astella/src/services/messageService"
	"github.com/kirtfieldk/astella/src/structures"
	"github.com/kirtfieldk/astella/src/util"
)

func (h *BaseHandler) PostMessageToEvent(c *gin.Context) {
	var msg structures.MessageRequestBody
	if err := c.BindJSON(&msg); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.PAYLOAD_IS_NOT_MSG})
		return
	}

	success, err := messageservice.PostMessage(msg, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{constants.MESSAGE: success})
}

func (h *BaseHandler) GetMessagesInEvent(c *gin.Context) {
	var point structures.Point
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	pagination := util.GetPageFromUrlQuery(c)
	if err := c.BindJSON(&point); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.PAYLOAD_IS_NOT_LOCATION})
		return
	}
	msg, err := messageservice.GetMessagesInEvent(eventId, userId, point, pagination, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, msg)
}

func (h *BaseHandler) UpvoteMessage(c *gin.Context) {
	var point structures.Point
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	messageId := c.Param(constants.MESSAGE_ID)
	if err := c.BindJSON(&point); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.PAYLOAD_IS_NOT_LOCATION})
		return
	}
	success, err := messageservice.UpvoteMessage(messageId, userId, eventId, point, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{constants.MESSAGE: success})
}

func (h *BaseHandler) GetUserUpvotes(c *gin.Context) {
	var point structures.Point
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	messageId := c.Param(constants.MESSAGE_ID)
	pagination := util.GetPageFromUrlQuery(c)

	if err := c.BindJSON(&point); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: constants.PAYLOAD_IS_NOT_LOCATION})
		return
	}
	users, err := messageservice.GetUserUpvotes(messageId, userId, eventId, point, pagination, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, users)
}
