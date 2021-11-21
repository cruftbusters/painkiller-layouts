package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type PendingRendersController struct {
	interval time.Duration
}

func (c *PendingRendersController) AddRoutes(router *httprouter.Router) {
	upgrader := websocket.Upgrader{}
	pendingRenders := make(chan types.Layout, 2)
	router.POST("/", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var layout types.Layout
		json.NewDecoder(r.Body).Decode(&layout)
		pendingRenders <- layout
		rw.WriteHeader(201)
	})
	router.GET("/", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		go func() {
			for {
				conn.WriteControl(websocket.PingMessage, nil, time.Time{})
				time.Sleep(c.interval)
			}
		}()
		if err := conn.WriteJSON(<-pendingRenders); err != nil {
			panic(err)
		}
	})
}
