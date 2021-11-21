package acceptance

import (
	"fmt"
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

	t.Run("distribute layouts", func(t *testing.T) {
		httpBaseURL, wsBaseURL := t2.TestServer(v1.Handler)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		channels := [2]chan types.Layout{}
		for i := 0; i < len(channels); i++ {
			channels[i] = make(chan types.Layout)
			conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
			t2.AssertNoError(t, err)
			defer conn.Close()

			go func(index int) {
				var layout types.Layout
				err := conn.ReadJSON(&layout)
				if err != nil {
					panic(err)
				}
				channels[index] <- layout
			}(i)
		}

		layouts := [2]types.Layout{}
		for i := 0; i < len(channels); i++ {
			layouts[i] = types.Layout{Id: fmt.Sprintf("layout #%d", i)}
			client.CreatePendingRender(t, layouts[i])
		}

		for i := 0; i < len(channels); i++ {
			select {
			case got := <-channels[i]:
				t2.AssertLayout(t, got, layouts[i])
			case <-time.After(time.Second):
				t.Error("expected notification in less than one second")
			}
		}
	})

	t.Run("buffer notifications", func(t *testing.T) {
		httpBaseURL, wsBaseURL := t2.TestServer(v1.Handler)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		layout := types.Layout{Id: "unhandled"}
		client.CreatePendingRender(t, layout)

		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn.Close()

		channel := make(chan types.Layout)
		go func() {
			var layout types.Layout
			err := conn.ReadJSON(&layout)
			if err != nil {
				panic(err)
			}
			channel <- layout
		}()

		select {
		case <-channel:
		case <-time.After(time.Second):
			t.Fatal("expected notification in less than one second")
		}
	})
}
