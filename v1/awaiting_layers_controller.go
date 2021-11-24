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
	awaitingHillshade AwaitingLayerService
}

func (c *AwaitingLayersController) AddRoutes(router *httprouter.Router) {
	router.POST("/v1/awaiting_heightmap", c.EnqueueLayout(c.awaitingHeightmap))
	router.POST("/v1/awaiting_hillshade", c.EnqueueLayout(c.awaitingHillshade))
	router.GET("/v1/awaiting_hillshade", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		AwaitingLayerServer(conn, c.awaitingHillshade)
		PingServer(conn)
	})
	router.GET("/v1/awaiting_heightmap", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		AwaitingLayerServer(conn, c.awaitingHeightmap)
		PingServer(conn)
	})
}

func PingServer(conn *websocket.Conn) {
	for {
		conn.WriteControl(websocket.PingMessage, nil, time.Time{})
		time.Sleep(5 * time.Second)
	}
}

func AwaitingLayerServer(conn *websocket.Conn, awaitingLayer AwaitingLayerService) {
	go func() {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		layout := awaitingLayer.Dequeue()
		conn.WriteJSON(layout)
		if _, _, err := conn.ReadMessage(); err != nil {
			awaitingLayer.Enqueue(layout)
			return
		}
	}()
}

func (c *AwaitingLayersController) EnqueueLayout(awaitingLayerService AwaitingLayerService) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var layout types.Layout
		if err := json.NewDecoder(r.Body).Decode(&layout); err != nil {
			rw.WriteHeader(400)
			return
		} else if err := awaitingLayerService.Enqueue(layout); err != nil {
			rw.WriteHeader(500)
			return
		}
		rw.WriteHeader(201)
	}
}
