package layouts

import (
	"bytes"
	"errors"
	"io"
	"testing"

	. "github.com/cruftbusters/painkiller-layouts/testing"
	"github.com/julienschmidt/httprouter"
)

func TestLayerController(t *testing.T) {
	mockLayerService := new(MockLayerService)
	controller := LayerController{
		mockLayerService,
	}
	client, _ := NewTestClient(t, func(string, string) *httprouter.Router {
		router := httprouter.New()
		controller.AddRoutes(router)
		return router
	})

	t.Run("put layer on missing map is not found", func(t *testing.T) {
		id, name := "there is no creativity", "heightmap.jpg"
		mockLayerService.On("Put", id, name, []byte{}).Return(ErrLayoutNotFound)
		client.PutLayerExpectNotFound(id, name)
	})

	t.Run("get layer on missing map is not found", func(t *testing.T) {
		id, name := "walrus", "heightmap.jpg"
		mockLayerService.On("Get", id, name).Return([]byte{}, "", ErrLayoutNotFound)
		client.GetLayerExpectNotFound(id, name)
	})

	t.Run("get layer is not found", func(t *testing.T) {
		id, name := "serendipity", "heightmap.jpg"
		mockLayerService.On("Get", id, name).Return([]byte{}, "", ErrLayerNotFound)
		client.GetLayerExpectNotFound(id, name)
	})

	t.Run("disallow put for invalid layer names", func(t *testing.T) {
		for _, name := range []string{
			"heightmap.tif",
			"hillshade.png",
			"anything.jpg",
		} {
			client.PutLayerExpectBadRequest("anything", name)
		}
	})

	t.Run("put layer", func(t *testing.T) {
		id, name, up := "john denver", "heightmap.jpg", []byte("was a bear")
		mockLayerService.On("Put", id, name, up).Return(nil)
		client.PutLayer(id, name, bytes.NewBuffer(up))
	})

	t.Run("get layer", func(t *testing.T) {
		id, name, layer, contentType := "inwards", "heightmap.jpg", []byte("buncha bytes"), "image/png"
		mockLayerService.On("Get", id, name).Return(layer, contentType, nil)

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
		mockLayerService.On("Delete", id, name).Return(nil)
		client.DeleteLayer(id, name)
	})

	t.Run("delete with error", func(t *testing.T) {
		id, name := "iron gear wheel", "hillshade.jpg"
		mockLayerService.On("Delete", id, name).Return(errors.New("anything"))
		client.DeleteLayerExpectInternalServerError(id, name)
	})
}
