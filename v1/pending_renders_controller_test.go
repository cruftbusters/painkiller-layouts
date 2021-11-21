package v1

import (
	"fmt"
	"sync"
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func TestPendingRendersController(t *testing.T) {
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

		layouts := [2]types.Layout{}
		for i := 0; i < len(layouts); i++ {
			layouts[i] = types.Layout{Id: fmt.Sprintf("layout #%d", i)}
			client.CreatePendingRender(t, layouts[i])
		}

		var wg sync.WaitGroup
		wg.Add(len(layouts))
		for _, want := range layouts {
			conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
			t2.AssertNoError(t, err)
			defer conn.Close()

			go func(want types.Layout) {
				got, err := (&t2.WSClient{Conn: conn}).ReadLayout()
				if err != nil {
					t.Errorf("got %s want nil", err)
				} else if got != want {
					t.Errorf("got %+v want %+v", got, want)
				}
				wg.Done()
			}(want)
		}
		wg.Wait()
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
		wsClient := t2.WSClient{Conn: conn}

		got, err := wsClient.ReadLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, layout)
	})

	t.Run("pull more work", func(t *testing.T) {
		controller := &PendingRendersController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn.Close()
		wsClient := t2.WSClient{Conn: conn}

		first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}

		client.CreatePendingRender(t, first)
		got, err := wsClient.ReadLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, first)

		client.CreatePendingRender(t, second)
		wsClient.Ready()
		got, err = wsClient.ReadLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
	})

	t.Run("redistribute abandoned work", func(t *testing.T) {
		controller := &PendingRendersController{time.Second}
		httpBaseURL, wsBaseURL := t2.TestController(controller)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		for i := 0; i < 16; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
			t2.AssertNoError(t, err)
			conn.Close()
		}

		conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn.Close()
		wsClient := t2.WSClient{Conn: conn}

		first, second := types.Layout{Id: "first"}, types.Layout{Id: "second"}

		client.CreatePendingRender(t, first)
		got, err := wsClient.ReadLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, first)

		client.CreatePendingRender(t, second)
		wsClient.Ready()
		got, err = wsClient.ReadLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
		conn.Close()

		conn, _, err = websocket.DefaultDialer.Dial(wsBaseURL, nil)
		t2.AssertNoError(t, err)
		defer conn.Close()
		wsClient = t2.WSClient{Conn: conn}
		got, err = wsClient.ReadLayout()
		t2.AssertNoError(t, err)
		t2.AssertLayout(t, got, second)
	})
}
