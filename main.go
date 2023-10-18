package main

import (
	"github.com/gin-gonic/gin"
	"gochatroom/api"
	"gochatroom/ws"
)

func main() {
	route := api.Route()
	ws.InitWs(route)

	gin.SetMode(gin.DebugMode)
	route.Run(":8080")
}
