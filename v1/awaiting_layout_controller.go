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
	interval time.Duration
}

func (c *AwaitingLayoutController) AddRoutes(router *httprouter.Router) {
	upgrader := websocket.Upgrader{}
	awaitingLayouts := make(chan types.Layout, 2)
	router.POST("/v1/layouts_awaiting", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var layout types.Layout
		json.NewDecoder(r.Body).Decode(&layout)
		select {
		case awaitingLayouts <- layout:
			rw.WriteHeader(201)
		default:
			log.Print("queue full")
			rw.WriteHeader(500)
		}
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

		for {
			awaitingLayout := <-awaitingLayouts
			if err := conn.WriteJSON(awaitingLayout); err != nil {
				log.Printf("failed websocket write: %s", err)
				awaitingLayouts <- awaitingLayout
				break
			} else if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("failed websocket read: %s", err)
				awaitingLayouts <- awaitingLayout
				break
			}
		}
	})
}
