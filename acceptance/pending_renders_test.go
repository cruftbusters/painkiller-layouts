package acceptance

import (
	"fmt"
	"sync"
	"testing"
	"time"

	t2 "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	v1 "github.com/cruftbusters/painkiller-layouts/v1"
	"github.com/gorilla/websocket"
)

func TestPendingRenders(t *testing.T) {
	t.Run("gracefully close connection", func(t *testing.T) {
		httpBaseURL, wsBaseURL := t2.TestServer(v1.Handler)
		client := t2.ClientV2{BaseURL: httpBaseURL}

		for i := 0; i < 16; i++ {
			conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL, nil)
			t2.AssertNoError(t, err)
			conn.Close()
		}

		client.CreatePendingRender(t, types.Layout{})
	})

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
		httpBaseURL, wsBaseURL := t2.TestServer(v1.Handler)
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
		httpBaseURL, wsBaseURL := t2.TestServer(v1.Handler)
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
}
