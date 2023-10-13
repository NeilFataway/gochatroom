package gochatroom

import (
	"gochatroom/api"
	"gochatroom/ws"
)

func main() {
	route := api.Route()
	ws.InitWs(route)

	route.Run(":8080")
}
