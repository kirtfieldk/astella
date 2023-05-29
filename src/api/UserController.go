package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kirtfieldk/astella/src/constants"
	s3service "github.com/kirtfieldk/astella/src/services/aws/S3Service"
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
	var username = c.Request.FormValue("username")
	fileNames := getFiles(c, h, username)
	user = structures.User{
		Id:          c.Request.FormValue("id"),
		Twitter:     c.Request.FormValue("twitter"),
		Ig:          c.Request.FormValue("ig"),
		TikTok:      c.Request.FormValue("tiktok"),
		Description: c.Request.FormValue("description"),
		ImgOne:      fileNames[0],
		ImgTwo:      fileNames[1],
		ImgThree:    fileNames[2],
		Username:    username,
	}
	log.Println(user)
	resp, err := userservice.UpdateUserProfile(user, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, resp)
}

func getFiles(c *gin.Context, h *BaseHandler, username string) [3]string {
	var resp [3]string
	var fileParams = [3]string{"img_one.png", "img_two.png", "img_three.png"}
	// c.Request.ParseMultipartForm(32 << 20)

	for i, el := range fileParams {
		fileHeader, err := c.FormFile(el)

		if err != nil {
			log.Println("error")
			log.Println(err)
		} else {
			var fileName = username + "_" + uuid.New().String()
			log.Printf("Created File with name %v\n", fileName)
			c.SaveUploadedFile(fileHeader, "tmp/"+fileName)

			fx, _ := os.Open("tmp/" + fileName)

			defer fx.Close()
			defer os.Remove(fx.Name())
			s3service.UploadObject("astellaapplicationmessages", fx, fileName, h.S3Session)
			resp[i] = fileName

		}

	}
	return resp

}

func (h *BaseHandler) GetUser(c *gin.Context) {
	userId := c.Param(constants.USER_ID)
	userResponse, err := userservice.GetUser(userId, h.DB)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: "No user"})
		return
	}
	c.IndentedJSON(http.StatusAccepted, userResponse)
}
