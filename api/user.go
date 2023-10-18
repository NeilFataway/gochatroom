package api

import (
	"github.com/gin-gonic/gin"
	"gochatroom/services"
)

type UserPo struct {
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
}

func CreateUser(c *gin.Context) {
	userPo := &UserPo{}
	if err := c.BindJSON(userPo); err != nil {
		services.WrapBadRequestError(err).ResponseError(c)
		return
	}

	if user, err := services.CreateUser(userPo.UserName, userPo.Avatar); err != nil {
		err.ResponseError(c)
		return
	} else {
		ResponseOK(c, user)
	}
}

func DeleteUser(c *gin.Context) {
	userId := c.Query("userId")
	if len(userId) == 0 {
		services.InValidUserID.ResponseError(c)
		return
	}

	if err := services.RemoveUser(userId); err != nil {
		err.ResponseError(c)
	} else {
		ResponseOK(c, "done")
	}
}

func GetUsers(c *gin.Context) {
	users := services.GetUsers()
	ResponseOK(c, users)
}
