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
				err := conn.WriteControl(websocket.PingMessage, nil, time.Time{})
				if err != nil {
					log.Printf("closing websocket: %s", err)
					break
				}
				time.Sleep(c.interval)
			}
		}()

		pendingRender := <-pendingRenders
		if err := conn.WriteJSON(pendingRender); err != nil {
			log.Printf("failed websocket write: %s", err)
		}
	})
}
