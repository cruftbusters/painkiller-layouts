package acceptance

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func TestLayers(t *testing.T) {
	httpBaseURL, _ := TestServer(layouts.Handler)
	client := ClientV2{BaseURL: httpBaseURL}

	id := client.CreateLayout(t, Layout{}).Id
	defer func() { client.DeleteLayout(t, id) }()

	t.Run("put layers on missing layout", func(t *testing.T) {
		client.PutLayerExpectNotFound(t, "deadbeef", "heightmap.jpg")
		client.PutLayerExpectNotFound(t, "deadbeef", "hillshade.jpg")
	})

	t.Run("get missing layers on layout", func(t *testing.T) {
		client.GetLayerExpectNotFound(t, id, "heightmap.jpg")
		client.GetLayerExpectNotFound(t, id, "hillshade.jpg")
	})

	for _, name := range []string{
		"heightmap.tif",
		"hillshade.png",
		"anything.jpg",
	} {
		t.Run(fmt.Sprintf("put '%s' bad request", name), func(t *testing.T) {
			client.PutLayerExpectBadRequest(t, id, name)
		})

		t.Run(fmt.Sprintf("get '%s' not found", name), func(t *testing.T) {
			client.GetLayerExpectNotFound(t, id, name)
		})
	}

	t.Run("put get delete", func(t *testing.T) {
		scenarios := []struct {
			name        string
			layer       []byte
			contentType string
		}{
			{
				name:        "heightmap.jpg",
				layer:       []byte("heightmap bytes"),
				contentType: "image/jpeg",
			},
			{
				name:        "hillshade.jpg",
				layer:       []byte("hillshade bytes"),
				contentType: "image/jpeg",
			},
		}

		for _, scenario := range scenarios {
			t.Run("put "+scenario.name, func(t *testing.T) {
				client.PutLayer(t, id, scenario.name, bytes.NewBuffer(scenario.layer))
			})
		}

		for _, scenario := range scenarios {
			t.Run("get "+scenario.name, func(t *testing.T) {
				gotReadCloser, gotContentType := client.GetLayer(t, id, scenario.name)
				got, err := io.ReadAll(gotReadCloser)
				AssertNoError(t, err)
				if !bytes.Equal(got, scenario.layer) {
					t.Errorf("got %v want %v", got, scenario.layer)
				}
				if gotContentType != scenario.contentType {
					t.Errorf("got %s want %s", gotContentType, scenario.contentType)
				}
			})
		}

		for _, scenario := range scenarios {
			t.Run("delete "+scenario.name, func(t *testing.T) {
				client.DeleteLayer(t, id, scenario.name)
				client.GetLayerExpectNotFound(t, id, scenario.name)
			})
		}
	})

	t.Run("put heightmap updates heightmap URL", func(t *testing.T) {
		client.PutLayer(t, id, "heightmap.jpg", nil)

		got := client.GetLayout(t, id).HeightmapURL
		want := fmt.Sprintf("%s/v1/layouts/%s/heightmap.jpg", httpBaseURL, id)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("put hillshade updates hillshade URL", func(t *testing.T) {
		client.PutLayer(t, id, "hillshade.jpg", nil)

		got := client.GetLayout(t, id).HillshadeURL
		want := fmt.Sprintf("%s/v1/layouts/%s/hillshade.jpg", httpBaseURL, id)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("layers are present after deleting layout", func(t *testing.T) {
		id := client.CreateLayout(t, Layout{}).Id
		client.PutLayer(t, id, "heightmap.jpg", nil)
		client.DeleteLayout(t, id)
		client.GetLayer(t, id, "heightmap.jpg")
	})
}
