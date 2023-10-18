package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"data":    data,
		"message": nil,
	})
}

func Route() *gin.Engine {
	r := gin.Default()
	apiGroup := r.Group("/api")
	initChatRoomRouter(apiGroup)
	initUserRouter(apiGroup)
	return r
}

func initChatRoomRouter(r *gin.RouterGroup) {
	r.POST("chat_rooms", CreateRoom)
	r.DELETE("chat_rooms", DeleteRoom)
	r.GET("chat_rooms", GetRooms)
	r.GET("chat_rooms/:room_id/users", GetRoomUsers)
}

func initUserRouter(r *gin.RouterGroup) {
	r.POST("users", CreateUser)
	r.DELETE("users", DeleteUser)
	r.GET("users", GetUsers)
}
