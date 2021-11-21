package acceptance

import (
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/gorilla/websocket"
)

func TestPendingRenders(t *testing.T) {
	t.Run("ping every interval", func(t *testing.T) {
		_, wsBaseURL := t2.TestServer(v1.Handler)
		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn.Close()

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
	})

	t.Run("broadcast layout", func(t *testing.T) {
		httpBaseURL, wsBaseURL := t2.TestServer(v1.Handler)
		client := t2.ClientV2{httpBaseURL}

		channel0 := make(chan types.Layout)
		conn0, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn0.Close()

		go func() {
			var layout types.Layout
			err := conn0.ReadJSON(&layout)
			if err != nil {
				panic(err)
			}
			channel0 <- layout
		}()

		channel1 := make(chan types.Layout)
		conn1, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn1.Close()

		go func() {
			var layout types.Layout
			err := conn1.ReadJSON(&layout)
			if err != nil {
				panic(err)
			}
			channel1 <- layout
		}()

		layout := types.Layout{Id: "notification"}
		client.CreatePendingRender(t, layout)

		select {
		case got := <-channel0:
			t2.AssertLayout(t, got, layout)
		case <-time.After(time.Second):
			t.Error("expected notification in less than one second")
		}

		select {
		case got := <-channel1:
			t2.AssertLayout(t, got, layout)
		case <-time.After(time.Second):
			t.Error("expected notification in less than one second")
		}
	})
}
