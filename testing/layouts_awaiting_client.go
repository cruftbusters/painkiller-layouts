package testing

import (
	"fmt"
	"time"

	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

type LayoutsAwaitingClient struct {
	Conn *websocket.Conn
}

func NewLayoutsAwaitingClient(wsBaseURL string) (LayoutsAwaitingClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/layouts_awaiting", nil)
	return LayoutsAwaitingClient{Conn: conn}, err
}

func (c *LayoutsAwaitingClient) StartDequeueAwaitingLayout() (types.Layout, error) {
	channel := make(chan struct {
		types.Layout
		error
	})
	go func() {
		var layout types.Layout
		err := c.Conn.ReadJSON(&layout)
		channel <- struct {
			types.Layout
			error
		}{layout, err}
	}()

	select {
	case result := <-channel:
		return result.Layout, result.error
	case <-time.After(time.Second):
		return types.Layout{}, fmt.Errorf("timed out after one second")
	}
}

func (c *LayoutsAwaitingClient) CompleteDequeueAwaitingLayout() error {
	return c.Conn.WriteMessage(websocket.BinaryMessage, nil)
}
