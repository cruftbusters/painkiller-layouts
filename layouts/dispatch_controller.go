package layouts

import (
	"encoding/json"
	"net/http"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type DispatchController struct {
	layoutPublisher chan types.Layout
}

func (c *DispatchController) AddRoutes(router *httprouter.Router) {
	upgrader := websocket.Upgrader{}
	var connection *websocket.Conn
	go func() {
		for {
			layout := <-c.layoutPublisher
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
		}
	}()
	router.GET("/v1/layout_dispatch", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var err error
		if connection, err = upgrader.Upgrade(rw, r, nil); err != nil {
			panic(err)
		}
	})
}
