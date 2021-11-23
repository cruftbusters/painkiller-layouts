package v1

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type AwaitingLayersController struct{}

func (c *AwaitingLayersController) AddRoutes(router *httprouter.Router) {
	router.GET("/v1/awaiting_heightmap", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		for {
			conn.WriteControl(websocket.PingMessage, nil, time.Time{})
			time.Sleep(5 * time.Second)
		}
	})
}
