package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
	var updatedImages = getFiles(c, h)
	user = structures.User{
		Id:          c.Request.FormValue("id"),
		Twitter:     c.Request.FormValue("twitter"),
		Ig:          c.Request.FormValue("ig"),
		TikTok:      c.Request.FormValue("tiktok"),
		Description: c.Request.FormValue("description"),
		ImgOne:      c.Request.FormValue("img_one"),
		ImgTwo:      c.Request.FormValue("img_two"),
		ImgThree:    c.Request.FormValue("img_three"),
		Username:    c.Request.FormValue("username"),
	}
	resp, err := userservice.UpdateUserProfile(user, updatedImages, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, resp)

}

func getFiles(c *gin.Context, h *BaseHandler) [3]string {
	var resp [3]string
	var fileParams = [3]string{"img_one", "img_two", "img_three"}
	c.Request.ParseMultipartForm(32 << 20)

	for i, el := range fileParams {
		fileHeader, err := c.FormFile(el)
		if err != nil {
			log.Println(err)
		} else {
			file, err := fileHeader.Open()
			if err != nil {
				log.Println(err)
			}
			defer file.Close()
			log.Println(fileHeader.Filename)
			log.Println(fileHeader.Size)
			tmpf, err := ioutil.TempFile("", fileHeader.Filename)
			if err != nil {
				log.Println(err)
			}
			defer tmpf.Close()
			defer os.Remove(tmpf.Name())

			path, err := s3service.UploadObject("astellaapplicationmessages", tmpf, fileHeader.Filename, h.AwsSession)
			resp[i] = path
		}

	}
	return resp

}
