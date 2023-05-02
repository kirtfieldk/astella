package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/constants"
	eventservices "github.com/kirtfieldk/astella/src/services/eventServices"
	"github.com/kirtfieldk/astella/src/structures"
	"github.com/kirtfieldk/astella/src/util"
)

type codeWithLocationPayload struct {
	Code      string  `json:"code"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
type cityPayload struct {
	City string `json:"city"`
}

func (h *BaseHandler) GetEvent(c *gin.Context) {
	id := c.Param(constants.ID)
	event, err := eventservices.GetEvent(id, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, event)

}

func (h *BaseHandler) GetEventByCity(c *gin.Context) {
	requestBody := readCityFromPayload(c)
	pagination := util.GetPageFromUrlQuery(c)
	events, err := eventservices.GetEventsByCity(requestBody.City, pagination, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err})
		return
	}
	c.IndentedJSON(http.StatusOK, events)
}

func (h *BaseHandler) CreateEvent(c *gin.Context) {
	var event structures.Event
	if err := c.BindJSON(&event); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: "Missing fields for event"})
		return
	}
	success, err := eventservices.CreateEvent(event, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{constants.MESSAGE: success})
}

func (h *BaseHandler) AddUserToEvent(c *gin.Context) {
	var requestBody codeWithLocationPayload
	eventId := c.Param(constants.EVENT_ID)
	userId := c.Param(constants.USER_ID)
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: constants.CODE_NOT_FOUND})
		return
	}
	point := structures.Point{Latitude: requestBody.Latitude, Longitude: requestBody.Longitude}
	success, err := eventservices.AddUserToEvent(requestBody.Code, userId, eventId, point, h.DB)
	if err != nil || !success {
		c.IndentedJSON(http.StatusNotFound, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{constants.MESSAGE: success})

}

func deleteEvent(c *gin.Context) {

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
