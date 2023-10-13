package api

import (
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	r := gin.Default()
	apiGroup := r.Group("/api")
	initChatRoomRouter(apiGroup)
	return r
}

func initChatRoomRouter(r *gin.RouterGroup) {
	r.POST("chatroom", CreateRoom)
	r.GET("chatroom", GetRooms)
}
