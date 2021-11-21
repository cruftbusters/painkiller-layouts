package acceptance

import (
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/gorilla/websocket"
)

func TestPendingRenders(t *testing.T) {
	_, wsBaseURL := t2.TestServer(v1.Handler)
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
