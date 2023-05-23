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
	messagesResp, err := messageservice.GetMessagesInEvent(msg.EventId, msg.UserId, 0, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	messagesResp.Success = success
	c.IndentedJSON(http.StatusAccepted, messagesResp)
}

func (h *BaseHandler) GetMessagesInEvent(c *gin.Context) {
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	pagination := util.GetPageFromUrlQuery(c)
	msg, err := messageservice.GetMessagesInEvent(eventId, userId, pagination, h.DB)
	if err != nil {
		print(err)
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
	message, err := messageservice.UpvoteMessage(messageId, userId, eventId, point, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, message)
}

func (h *BaseHandler) FetchMessageThread(c *gin.Context) {
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	messageId := c.Param(constants.MESSAGE_ID)
	page := util.GetPageFromUrlQuery(c)
	resp, err := messageservice.GetMessageThread(messageId, userId, eventId, page, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, resp)
}

func (h *BaseHandler) DownvoteMessage(c *gin.Context) {
	var point structures.Point
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	messageId := c.Param(constants.MESSAGE_ID)
	if err := c.BindJSON(&point); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.PAYLOAD_IS_NOT_LOCATION})
		return
	}
	message, err := messageservice.DownVoteMessage(messageId, userId, eventId, point, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, message)
}

func (h *BaseHandler) GetUserUpvotes(c *gin.Context) {
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	messageId := c.Param(constants.MESSAGE_ID)
	pagination := util.GetPageFromUrlQuery(c)
	users, err := messageservice.GetUserUpvotes(messageId, userId, eventId, pagination, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, users)
}

func (h *BaseHandler) PinMessage(c *gin.Context) {
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	messageId := c.Param(constants.MESSAGE_ID)
	success, err := messageservice.PinnMessage(messageId, userId, eventId, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, success)
}

func (h *BaseHandler) UnpinMessage(c *gin.Context) {
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	messageId := c.Param(constants.MESSAGE_ID)
	success, err := messageservice.UnpinnMessage(messageId, userId, eventId, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, success)
}

func (h *BaseHandler) GetPinnedMessaged(c *gin.Context) {
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	pagination := util.GetPageFromUrlQuery(c)
	messages, err := messageservice.GetUsersPinnedMessagesInEvent(userId, eventId, pagination, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, messages)
}
