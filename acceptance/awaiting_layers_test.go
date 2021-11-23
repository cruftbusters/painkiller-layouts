package acceptance

import (
	"testing"
	"time"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/gorilla/websocket"
)

func TestAwaitingLayers(t *testing.T) {
	_, wsBaseURL := TestServer(v1.Handler)

	t.Run("ping every five seconds", func(t *testing.T) {
		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/awaiting_heightmap", nil)
		AssertNoError(t, err)
		defer conn.Close()
		go conn.ReadMessage()

		ping := make(chan *struct{})
		conn.SetPingHandler(func(string) error { ping <- nil; return nil })

		one, five, six := time.After(time.Second), time.After(5*time.Second), time.After(6*time.Second)
		select {
		case <-ping:
		case <-one:
			t.Fatal("timed out waiting for first ping")
		}

		select {
		case <-ping:
			t.Fatal("second ping too early")
		case <-five:
		}

		select {
		case <-ping:
		case <-six:
			t.Fatal("second ping too late")
		}
	})
}
