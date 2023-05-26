package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/api/responses"
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
	fileNames := getFiles(c, h)
	user = structures.User{
		Id:          c.Request.FormValue("id"),
		Twitter:     c.Request.FormValue("twitter"),
		Ig:          c.Request.FormValue("ig"),
		TikTok:      c.Request.FormValue("tiktok"),
		Description: c.Request.FormValue("description"),
		ImgOne:      fileNames[0],
		ImgTwo:      fileNames[1],
		ImgThree:    fileNames[2],
		Username:    c.Request.FormValue("username"),
	}
	resp, err := userservice.UpdateUserProfile(user, h.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{constants.MESSAGE: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, resp)
}

func getFiles(c *gin.Context, h *BaseHandler) [3]string {
	var resp [3]string
	var fileParams = [3]string{"img_one", "img_two", "img_three"}
	// c.Request.ParseMultipartForm(32 << 20)

	for i, el := range fileParams {
		fileHeader, err := c.FormFile(el)

		if err != nil {
			log.Println(err)
		} else {
			log.Println(fileHeader.Filename)
			c.SaveUploadedFile(fileHeader, "tmp/"+fileHeader.Filename)

			fx, err := os.Open("tmp/" + fileHeader.Filename)
			if err != nil {
				log.Println("FUCK ")
			}
			defer fx.Close()
			defer os.Remove(fx.Name())
			s3service.UploadObject("astellaapplicationmessages", fx, fileHeader.Filename, h.S3Session)
			resp[i] = fileHeader.Filename

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
	h.getPresignedUrls(userResponse)
	c.IndentedJSON(http.StatusAccepted, userResponse)

}

func (h *BaseHandler) getPresignedUrls(usr responses.UserListResponse) {
	var result []string
	if len(usr.Data) > 0 {
		images := []string{usr.Data[0].ImgOne, usr.Data[0].ImgTwo, usr.Data[0].ImgThree}
		for _, el := range images {
			log.Println(el)
			url, err := s3service.GeneratePresignedUrl(h.AwsBucket, el, h.S3Session)
			if err != nil {
				log.Println(err)
				continue
			}
			result = append(result, url)
		}
		usr.Data[0].PresignedUrl = result
	}
}
