package layouts

import (
	"fmt"
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func TestDispatch(t *testing.T) {
	layoutPublisher := make(chan types.Layout)
	pingInterval := time.Second
	controller := &DispatchController{layoutPublisher, pingInterval}

	_, wsBaseURL := TestController(controller)

	t.Run("sink layouts when no subscribers", func(t *testing.T) {
		down := types.Layout{Id: "hello im new here"}
		layoutPublisher <- down
	})

	t.Run("server ping", func(t *testing.T) {
		connection, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/v1/layout_dispatch", wsBaseURL), nil)
		AssertNoError(t, err)
		defer connection.Close()
		go func() {
			for {
				_, reader, err := connection.NextReader()
				if err != nil {
					break
				}
				reader.Read([]byte{})
			}
		}()

		signal := make(chan *struct{})
		connection.SetPingHandler(func(string) error { signal <- nil; return nil })
		select {
		case <-signal:
		case <-time.After(pingInterval + time.Second):
			t.Fatalf("No ping after %s", pingInterval+time.Second)
		}

		select {
		case <-signal:
			t.Fatalf("Next ping did not wait at least %s", pingInterval)
		case <-time.After(pingInterval):
		}
	})

	t.Run("dispatch new layouts", func(t *testing.T) {
		connection, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/v1/layout_dispatch", wsBaseURL), nil)
		AssertNoError(t, err)
		defer connection.Close()

		down := types.Layout{Id: "hello im new here"}
		layoutPublisher <- down

		var layout types.Layout
		err = connection.ReadJSON(&layout)
		AssertNoError(t, err)
		AssertLayout(t, layout, down)
	})
}
