package layouts

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/cruftbusters/painkiller-layouts/types"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

func TestDispatch(t *testing.T) {
	listener, protolessBaseURL := RandomPortListener()
	layoutPublisher := make(chan types.Layout)
	router := httprouter.New()
	(&DispatchController{layoutPublisher}).AddRoutes(router)
	go func() { http.Serve(listener, router) }()

	t.Run("sink layouts when no subscribers", func(t *testing.T) {
		down := types.Layout{Id: "hello im new here"}
		layoutPublisher <- down
	})

	connection, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/v1/layout_dispatch", "ws://"+protolessBaseURL), nil)
	AssertNoError(t, err)

	t.Run("dispatch new layouts", func(t *testing.T) {
		down := types.Layout{Id: "hello im new here"}
		layoutPublisher <- down

		var layout types.Layout
		err := connection.ReadJSON(&layout)
		AssertNoError(t, err)
		AssertLayout(t, layout, down)
	})
}
