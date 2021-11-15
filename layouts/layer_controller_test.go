package layouts

import (
	"bytes"
	"io"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/julienschmidt/httprouter"
)

func TestLayerController(t *testing.T) {
	stubLayerService := &StubLayerService{t: t}
	controller := LayerController{
		stubLayerService,
	}
	client, _ := NewTestClient(t, func(string, string) *httprouter.Router {
		router := httprouter.New()
		controller.AddRoutes(router)
		return router
	})

	t.Run("put layer on missing map is not found", func(t *testing.T) {
		id, name := "there is no creativity", "heightmap.jpg"
		stubLayerService.whenPutCalledWithId = id
		stubLayerService.putWillReturn = ErrLayoutNotFound

		client.PutLayerExpectNotFound(id, name)
	})

	t.Run("get layer on missing map is not found", func(t *testing.T) {
		id, name := "walrus", "heightmap.jpg"
		stubLayerService.whenGetCalledWithId = id
		stubLayerService.getWillReturnError = ErrLayoutNotFound

		client.GetLayerExpectNotFound(id, name)
	})

	t.Run("get layer is not found", func(t *testing.T) {
		id, name := "serendipity", "heightmap.jpg"
		stubLayerService.whenGetCalledWithId = id
		stubLayerService.getWillReturnError = ErrLayerNotFound

		client.GetLayerExpectNotFound(id, name)
	})

	t.Run("put layer", func(t *testing.T) {
		id, name, up := "john denver", "heightmap.jpg", []byte("was a bear")
		stubLayerService.whenPutCalledWithId = id
		stubLayerService.whenPutCalledWithLayer = up
		stubLayerService.putWillReturn = nil

		client.PutLayer(id, name, bytes.NewBuffer(up))
	})

	t.Run("get layer", func(t *testing.T) {
		id, name, layer, contentType := "inwards", "heightmap.jpg", []byte("buncha bytes"), "image/png"
		stubLayerService.whenGetCalledWithId = id
		stubLayerService.getWillReturnLayer = layer
		stubLayerService.getWillReturnContentType = contentType
		stubLayerService.getWillReturnError = nil

		gotReadCloser, gotContentType := client.GetLayer(id, name)
		got, err := io.ReadAll(gotReadCloser)
		AssertNoError(t, err)
		want := layer
		if !bytes.Equal(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
		if gotContentType != contentType {
			t.Errorf("got %s want %s", gotContentType, contentType)
		}
	})
}
