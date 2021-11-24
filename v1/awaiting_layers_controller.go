package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type AwaitingLayersController struct {
	awaitingHeightmap AwaitingLayerService
}

func (c *AwaitingLayersController) AddRoutes(router *httprouter.Router) {
	router.POST("/v1/awaiting_heightmap", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var layout types.Layout
		if err := json.NewDecoder(r.Body).Decode(&layout); err != nil {
			rw.WriteHeader(400)
			return
		} else if err := c.awaitingHeightmap.Enqueue(layout); err != nil {
			rw.WriteHeader(500)
			return
		}
		rw.WriteHeader(201)
	})
	router.GET("/v1/awaiting_heightmap", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		go func() {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
			layout := c.awaitingHeightmap.Dequeue()
			conn.WriteJSON(layout)
			if _, _, err := conn.ReadMessage(); err != nil {
				c.awaitingHeightmap.Enqueue(layout)
				return
			}
		}()
		for {
			conn.WriteControl(websocket.PingMessage, nil, time.Time{})
			time.Sleep(5 * time.Second)
		}
	})
}
