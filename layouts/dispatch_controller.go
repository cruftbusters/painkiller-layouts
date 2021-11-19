package layouts

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type DispatchController struct {
	layoutPublisher chan types.Layout
	pingInterval    time.Duration
}

func (c *DispatchController) AddRoutes(router *httprouter.Router) {
	upgrader := websocket.Upgrader{}
	var connection *websocket.Conn
	go func() {
		for {
			select {
			case layout := <-c.layoutPublisher:
				if connection != nil {
					writer, err := connection.NextWriter(websocket.TextMessage)
					if err != nil {
						panic(err)
					}
					if err := json.NewEncoder(writer).Encode(layout); err != nil {
						panic(err)
					}
					if err := writer.Close(); err != nil {
						panic(err)
					}
				}
			case <-time.After(c.pingInterval):
				if connection != nil {
					err := connection.WriteControl(websocket.PingMessage, nil, time.Time{})
					if err != nil {
						connection = nil
					}
				}
			}
		}
	}()
	router.GET("/v1/layout_dispatch", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var err error
		if connection, err = upgrader.Upgrade(rw, r, nil); err != nil {
			panic(err)
		}
		connection.SetCloseHandler(func(int, string) error { connection = nil; return nil })
	})
}
