package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
	"github.com/kirtfieldk/astella/src/services"
)

type codePayload struct {
	Code string `json:"code"`
}
type cityPayload struct {
	City string `json:"city"`
}

func GetEvent(c *gin.Context) {
	id := c.Param(constants.ID)
	event, err := services.GetEvent(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err})
		return
	}
	c.IndentedJSON(http.StatusOK, event)

}

func GetEventByCity(c *gin.Context) {
	requestBody := readCityFromPayload(c)
	events, err := services.GetEventsByCity(requestBody.City)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err})
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

func AddUserToEvent(c *gin.Context) {
	requestBody := readCodeFromPayload(c)
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.CODE_NOT_FOUND})
		return
	}
	success, err := services.AddUserToEvent(requestBody.Code, userId, eventId)
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
