package v1

import (
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/gorilla/websocket"
)

func TestPendingRendersController(t *testing.T) {
	controller := &PendingRendersController{}
	_, wsBaseURL := t2.TestController(controller)
	conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
	t2.AssertNoError(t, err)

	go conn.ReadMessage()

	signal := make(chan *struct{})
	conn.SetPingHandler(func(string) error { signal <- nil; return nil })

	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Error("expected ping in less than one second")
	}
}
