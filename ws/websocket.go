package ws

import (
	"github.com/gin-gonic/gin"
)

func InitWs(router *gin.Engine) {
	router.GET("/ws", chat)
}
