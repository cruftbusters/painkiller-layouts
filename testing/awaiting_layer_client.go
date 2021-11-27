package testing

import (
	"errors"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func BeginDequeueLayout(conn *websocket.Conn, priority int) (types.Layout, error) {
	if err := conn.WriteMessage(websocket.BinaryMessage, nil); err != nil {
		return types.Layout{}, err
	}
	channel := make(chan struct {
		types.Layout
		error
	})
	go func() {
		var got types.Layout
		err := conn.ReadJSON(&got)
		channel <- struct {
			types.Layout
			error
		}{got, err}
	}()
	select {
	case result := <-channel:
		return result.Layout, result.error
	case <-time.After(time.Second):
		return types.Layout{}, errors.New("timed out after one second")
	}
}

func EndDequeueLayout(conn *websocket.Conn) error {
	return conn.WriteMessage(websocket.BinaryMessage, nil)
}
