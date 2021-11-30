package testing

import (
	"errors"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func BeginDequeueLayout(conn *websocket.Conn) (types.Layout, error) {
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

func DrainAwaitingLayers(wsBaseURL string) (func(), error) {
	c0, err := DrainAwaitingLayer(wsBaseURL + "/v1/awaiting_heightmap")
	if err != nil {
		return nil, err
	}
	c1, err := DrainAwaitingLayer(wsBaseURL + "/v1/awaiting_hillshade")
	if err != nil {
		c0.Close()
		return nil, err
	}
	return func() {
		c0.Close()
		c1.Close()
	}, nil
}

func DrainAwaitingLayer(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return conn, err
	}
	go func() {
		for {
			if _, err := BeginDequeueLayout(conn); err != nil {
				return
			} else if err := EndDequeueLayout(conn); err != nil {
				return
			}
		}
	}()
	return conn, nil
}
