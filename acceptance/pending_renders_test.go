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

	one, five, six := time.After(time.Second), time.After(5*time.Second), time.After(6*time.Second)
	select {
	case <-signal:
	case <-one:
		t.Fatal("expected ping in less than one second")
	}

	select {
	case <-signal:
		t.Fatal("expected no ping in less than five seconds")
	case <-five:
	}

	select {
	case <-signal:
	case <-six:
		t.Fatal("expected ping in less than six seconds")
	}
}
