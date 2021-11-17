package layouts

import (
	"bytes"
	"errors"
	"fmt"
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
		stubLayerService.whenPutCalledWithName = name
		stubLayerService.putWillReturn = ErrLayoutNotFound

		client.PutLayerExpectNotFound(id, name)
	})

	t.Run("get layer on missing map is not found", func(t *testing.T) {
		id, name := "walrus", "heightmap.jpg"
		stubLayerService.whenGetCalledWithId = id
		stubLayerService.whenGetCalledWithName = name
		stubLayerService.getWillReturnError = ErrLayoutNotFound

		client.GetLayerExpectNotFound(id, name)
	})

	t.Run("get layer is not found", func(t *testing.T) {
		id, name := "serendipity", "heightmap.jpg"
		stubLayerService.whenGetCalledWithId = id
		stubLayerService.whenGetCalledWithName = name
		stubLayerService.getWillReturnError = ErrLayerNotFound

		client.GetLayerExpectNotFound(id, name)
	})

	t.Run("constrain put and get", func(t *testing.T) {
		for _, name := range []string{
			"heightmap.tif",
			"hillshade.png",
			"anything.jpg",
		} {
			t.Run(fmt.Sprintf("disallow put '%s'", name), func(t *testing.T) {
				client.PutLayerExpectBadRequest("anything", name)
			})

			t.Run(fmt.Sprintf("'%s' not found", name), func(t *testing.T) {
				client.GetLayerExpectNotFound("anything", name)
			})
		}
	})

	t.Run("put layer", func(t *testing.T) {
		id, name, up := "john denver", "heightmap.jpg", []byte("was a bear")
		stubLayerService.whenPutCalledWithId = id
		stubLayerService.whenPutCalledWithName = name
		stubLayerService.whenPutCalledWithLayer = up
		stubLayerService.putWillReturn = nil

		client.PutLayer(id, name, bytes.NewBuffer(up))
	})

	t.Run("get layer", func(t *testing.T) {
		id, name, layer, contentType := "inwards", "heightmap.jpg", []byte("buncha bytes"), "image/png"
		stubLayerService.whenGetCalledWithId = id
		stubLayerService.whenGetCalledWithName = name
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

	t.Run("delete", func(t *testing.T) {
		id, name := "floral", "heightmap.jpg"

		stubLayerService.whenDeleteCalledWithId = id
		stubLayerService.whenDeleteCalledWithName = name
		stubLayerService.deleteWillReturn = nil

		client.DeleteLayer(id, name)
	})

	t.Run("delete with error", func(t *testing.T) {
		id, name := "iron gear wheel", "hillshade.jpg"

		stubLayerService.whenDeleteCalledWithId = id
		stubLayerService.whenDeleteCalledWithName = name
		stubLayerService.deleteWillReturn = errors.New("anything")

		client.DeleteLayerExpectInternalServerError(id, name)
	})
}
