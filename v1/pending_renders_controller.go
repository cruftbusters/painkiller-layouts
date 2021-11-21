package v1

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type PendingRendersController struct{}

func (c *PendingRendersController) AddRoutes(router *httprouter.Router) {
	upgrader := websocket.Upgrader{}
	router.GET("/", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		conn.WriteControl(websocket.PingMessage, nil, time.Time{})
	})
}
