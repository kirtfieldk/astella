package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
	userservice "github.com/kirtfieldk/astella/src/services/userService"
	"github.com/kirtfieldk/astella/src/structures"
	"github.com/kirtfieldk/astella/src/util"
)

func (h *BaseHandler) GeteventsMemberOf(c *gin.Context) {
	userId := c.Param(constants.USER_ID)
	pagination := util.GetPageFromUrlQuery(c)
	events, err := userservice.GetEventUserIsMember(userId, pagination, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, events)

}

func (h *BaseHandler) GeteventsMembers(c *gin.Context) {
	eventId := c.Param(constants.EVENT_ID)
	pagination := util.GetPageFromUrlQuery(c)
	events, err := userservice.GetEventMembers(eventId, pagination, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, events)

}

func (h *BaseHandler) UpdateUser(c *gin.Context) {
	var user structures.User
	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: "Missing fields for user"})
		return
	}
	resp, err := userservice.UpdateUserProfile(user, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, resp)

}
