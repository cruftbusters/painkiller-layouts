package v1

import (
	"encoding/json"
	"net/http"
	"strconv"
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
		AwaitingLayerServer(conn, c.awaitingHillshade, priority(r))
		PingServer(conn)
	})
	router.GET("/v1/awaiting_heightmap", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		AwaitingLayerServer(conn, c.awaitingHeightmap, priority(r))
		PingServer(conn)
	})
}

func priority(r *http.Request) int {
	priority, err := strconv.Atoi(r.URL.Query().Get("priority"))
	if err != nil {
		return 0
	}
	return priority
}

func PingServer(conn *websocket.Conn) {
	conn.SetPingHandler(func(s string) error { return conn.WriteControl(websocket.PongMessage, []byte(s), time.Time{}) })
	for {
		conn.WriteControl(websocket.PingMessage, nil, time.Time{})
		time.Sleep(5 * time.Second)
	}
}

func AwaitingLayerServer(conn *websocket.Conn, awaitingLayer AwaitingLayerService, priority int) {
	read := make(chan error, 1)
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			read <- err
			if err != nil {
				return
			}
		}
	}()
	go func() {
		for {
			if err := <-read; err != nil {
				return
			}
			layout := awaitingLayer.Dequeue(priority)
			if err := conn.WriteJSON(layout); err != nil {
				awaitingLayer.Enqueue(layout)
				return
			} else if err := <-read; err != nil {
				awaitingLayer.Enqueue(layout)
				return
			}
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
