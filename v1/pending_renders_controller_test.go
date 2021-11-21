package v1

import (
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func TestPendingRendersController(t *testing.T) {
	t.Run("ping every interval", func(t *testing.T) {
		interval := time.Second
		controller := &PendingRendersController{interval}
		_, wsBaseURL := t2.TestController(controller)
		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)

		go conn.ReadMessage()

		signal := make(chan *struct{})
		conn.SetPingHandler(func(string) error { signal <- nil; return nil })

		one, five, six := time.After(time.Second), time.After(interval), time.After(interval+time.Second)
		select {
		case <-signal:
		case <-one:
			t.Fatal("expected ping in less than one second")
		}

		select {
		case <-signal:
			t.Fatalf("expected no ping in less than %s", interval)
		case <-five:
		}

		select {
		case <-signal:
		case <-six:
			t.Fatalf("expected ping in less than %s", interval+time.Second)
		}
	})

	t.Run("broadcast layout", func(t *testing.T) {
		controller := &PendingRendersController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}
		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn.Close()

		channel := make(chan types.Layout)
		go func() {
			var layout types.Layout
			err := conn.ReadJSON(&layout)
			t2.AssertNoError(t, err)
			channel <- layout
		}()

		layout := types.Layout{Id: "notification"}
		client.CreatePendingRender(t, layout)

		select {
		case got := <-channel:
			t2.AssertLayout(t, got, layout)
		case <-time.After(time.Second):
			t.Error("expected notification in less than one second")
		}
	})
}
