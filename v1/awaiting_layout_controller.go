package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type AwaitingLayoutController struct {
	interval                 time.Duration
	layoutsAwaitingHeightmap chan types.Layout
	layoutsAwaitingHillshade chan types.Layout
}

func (c *AwaitingLayoutController) AddRoutes(router *httprouter.Router) {
	upgrader := websocket.Upgrader{}
	router.POST("/v1/layouts_awaiting", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var layout types.Layout
		json.NewDecoder(r.Body).Decode(&layout)
		select {
		case c.layoutsAwaitingHeightmap <- layout:
			rw.WriteHeader(201)
		default:
			log.Print("queue full")
			rw.WriteHeader(500)
		}
	})
	router.GET("/v1/layouts_awaiting", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		go func() {
			for {
				err := conn.WriteControl(websocket.PingMessage, nil, time.Time{})
				if err != nil {
					log.Printf("closing websocket: %s", err)
					break
				}
				time.Sleep(c.interval)
			}
		}()

		var channel chan types.Layout
		if r.URL.Query().Get("layer") == "hillshade" {
			channel = c.layoutsAwaitingHillshade
		} else {
			channel = c.layoutsAwaitingHeightmap
		}
		for {
			layout := <-channel
			if err := conn.WriteJSON(layout); err != nil {
				log.Printf("failed websocket write: %s", err)
				channel <- layout
				break
			} else if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("failed websocket read: %s", err)
				channel <- layout
				break
			}
		}
	})
}
