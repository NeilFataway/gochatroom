package api

import (
	"github.com/gin-gonic/gin"
	"gochatroom/services"
)

type RoomPo struct {
	RoomName string `json:"room_name"`
	OwnerId  string `json:"owner_id"`
}

func CreateRoom(c *gin.Context) {
	// 解析body
	roomPo := &RoomPo{}
	if err := c.BindJSON(roomPo); err != nil {
		services.WrapBadRequestError(err).ResponseError(c)
		return
	}

	// 获取owner用户
	owner, err := services.GetUser(roomPo.OwnerId)
	if err != nil {
		services.WrapBadRequestError(err).ResponseError(c)
		return
	}

	// 创建房间
	room, err := services.CreateRoom(roomPo.RoomName, owner)
	if err != nil {
		err.ResponseError(c)
	} else {
		ResponseOK(c, room)
	}
}

func GetRooms(c *gin.Context) {
	rooms := services.GetRooms()
	ResponseOK(c, rooms)
}

func DeleteRoom(c *gin.Context) {
	userId := c.Query("userId")
	if len(userId) == 0 {
		services.InValidUserID.ResponseError(c)
		return
	}

	roomId := c.Query("roomId")
	if len(roomId) == 0 {
		services.InValidRoomName.ResponseError(c)
		return
	}

	if err := services.DeleteRoom(roomId, userId); err != nil {
		err.ResponseError(c)
		return
	}

	ResponseOK(c, "done")
}

func GetRoomUsers(c *gin.Context) {
	roomId := c.Param("room_id")
	if len(roomId) == 0 {
		services.InValidRoomName.ResponseError(c)
		return
	}

	if users, err := services.GetRoomUsers(roomId); err != nil {
		err.ResponseError(c)
		return
	} else {
		ResponseOK(c, users)
	}
}
