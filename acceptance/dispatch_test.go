package acceptance

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func TestDispatch(t *testing.T) {
	listener, httpBaseURL, wsBaseURL := TestServer()
	router := layouts.Handler("file::memory:?cache=shared", httpBaseURL)
	go func() { http.Serve(listener, router) }()

	client := ClientV2{BaseURL: httpBaseURL}

	t.Run("server ping", func(t *testing.T) {
		t.SkipNow()

		connection, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/layout_dispatch", nil)
		AssertNoError(t, err)

		signal := make(chan *struct{})
		connection.SetPingHandler(func(string) error { signal <- nil; return nil })

		go func() {
			for {
				_, reader, err := connection.NextReader()
				if err != nil {
					break
				}
				reader.Read([]byte{})
			}
		}()

		select {
		case <-signal:
		case <-time.After(4 * time.Second):
			t.Fatal("No ping after 4 seconds")
		}

		select {
		case <-signal:
			t.Fatal("Next ping did not wait at least 3 seconds")
		case <-time.After(3 * time.Second):
		}
	})

	t.Run("dispatch new layouts", func(t *testing.T) {
		connection, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/v1/layout_dispatch", nil)
		AssertNoError(t, err)

		var wg sync.WaitGroup
		wg.Add(2)
		var gotDispatched, gotCreated types.Layout

		go func() {
			err := connection.ReadJSON(&gotDispatched)
			AssertNoError(t, err)
			wg.Done()
		}()

		go func() {
			gotCreated = client.CreateLayout(t, types.Layout{})
			client.DeleteLayout(t, gotCreated.Id)
			wg.Done()
		}()

		wg.Wait()

		AssertLayout(t, gotDispatched, gotCreated)
	})
}
