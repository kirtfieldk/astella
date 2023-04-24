package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
	eventservices "github.com/kirtfieldk/astella/src/services/eventServices"
)

type codePayload struct {
	Code string `json:"code"`
}
type cityPayload struct {
	City string `json:"city"`
}

func (h *BaseHandler) GetEvent(c *gin.Context) {
	id := c.Param(constants.ID)
	event, err := eventservices.GetEvent(id, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err})
		return
	}
	c.IndentedJSON(http.StatusOK, event)

}

func (h *BaseHandler) GetEventByCity(c *gin.Context) {
	requestBody := readCityFromPayload(c)
	events, err := eventservices.GetEventsByCity(requestBody.City, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err})
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

func (h *BaseHandler) AddUserToEvent(c *gin.Context) {
	requestBody := readCodeFromPayload(c)
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.CODE_NOT_FOUND})
		return
	}
	success, err := eventservices.AddUserToEvent(requestBody.Code, userId, eventId, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.NO_EVENT_FOUND})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{constants.MESSAGE: success})

}

func deleteEvent(c *gin.Context) {

}

func readCodeFromPayload(c *gin.Context) *codePayload {

	var requestBody codePayload
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.CODE_NOT_FOUND})
		return nil
	}
	return &requestBody
}

func readCityFromPayload(c *gin.Context) *cityPayload {
	var requestBody cityPayload
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.CITY_NOT_FOUND})
		return nil
	}
	return &requestBody
}
