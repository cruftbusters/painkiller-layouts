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
}
