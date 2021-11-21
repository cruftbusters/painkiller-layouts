package v1

import (
	"fmt"
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

		channels := [2]chan struct {
			types.Layout
			error
		}{}
		for i := 0; i < len(channels); i++ {
			channels[i] = make(chan struct {
				types.Layout
				error
			})
			conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
			t2.AssertNoError(t, err)
			defer conn.Close()

			go func(index int) {
				var layout types.Layout
				err := conn.ReadJSON(&layout)
				channels[index] <- struct {
					types.Layout
					error
				}{layout, err}
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
				t2.AssertNoError(t, got.error)
				t2.AssertLayout(t, got.Layout, layouts[i])
			case <-time.After(time.Second):
				t.Error("expected notification in less than one second")
			}
		}
	})

	t.Run("buffer notifications", func(t *testing.T) {
		controller := &PendingRendersController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		layout := types.Layout{Id: "unhandled"}
		client.CreatePendingRender(t, layout)

		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn.Close()

		channel := make(chan struct {
			types.Layout
			error
		})
		go func() {
			var layout types.Layout
			err := conn.ReadJSON(&layout)
			channel <- struct {
				types.Layout
				error
			}{layout, err}
		}()

		select {
		case result := <-channel:
			t2.AssertNoError(t, result.error)
			t2.AssertLayout(t, result.Layout, layout)
		case <-time.After(time.Second):
			t.Fatal("expected notification in less than one second")
		}
	})

	t.Run("gracefully close connection", func(t *testing.T) {
		controller := &PendingRendersController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		for i := 0; i < 16; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
			t2.AssertNoError(t, err)
			conn.Close()
		}

		client.CreatePendingRender(t, types.Layout{})
	})
}
