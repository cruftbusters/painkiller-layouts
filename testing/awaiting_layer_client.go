package testing

import (
	"errors"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func ReadLayout(conn *websocket.Conn) (types.Layout, error) {
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
