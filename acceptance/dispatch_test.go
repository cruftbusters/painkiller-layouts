package acceptance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
)

func TestDispatch(t *testing.T) {
	t.Run("dispatch new layouts", func(t *testing.T) {
		listener, protolessBaseURL := RandomPortListener()
		router := layouts.Handler("file::memory:?cache=shared", "http://"+protolessBaseURL)
		go func() { http.Serve(listener, router) }()

		client := ClientV2{BaseURL: "http://" + protolessBaseURL}
		connection, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/v1/layout_dispatch", "ws://"+protolessBaseURL), nil)
		AssertNoError(t, err)

		var wg sync.WaitGroup
		wg.Add(2)
		var gotDispatched, gotCreated types.Layout

		go func() {
			messageType, reader, err := connection.NextReader()
			AssertNoError(t, err)
			wantMessageType := websocket.TextMessage
			if messageType != wantMessageType {
				t.Errorf("got %d want %d", messageType, wantMessageType)
			}

			json.NewDecoder(reader).Decode(&gotDispatched)
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
