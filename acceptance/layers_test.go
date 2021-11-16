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
	client, baseURL := NewTestClient(t, layouts.Handler)

	t.Run("put layers on missing layout", func(t *testing.T) {
		client.PutLayerExpectNotFound("deadbeef", "heightmap.jpg")
		client.PutLayerExpectNotFound("deadbeef", "hillshade.jpg")
	})

	t.Run("get layers on missing layout", func(t *testing.T) {
		client.GetLayerExpectNotFound("deadbeef", "heightmap.jpg")
		client.GetLayerExpectNotFound("deadbeef", "hillshade.jpg")
	})

	t.Run("get missing layers on layout", func(t *testing.T) {
		id := client.CreateLayout(Layout{}).Id
		defer func() { client.DeleteLayout(id) }()
		client.GetLayerExpectNotFound(id, "heightmap.jpg")
		client.GetLayerExpectNotFound(id, "hillshade.jpg")
	})

	for _, scenario := range []struct {
		name        string
		layer       []byte
		contentType string
	}{
		{
			name:        "heightmap.jpg",
			layer:       []byte{65, 66, 67},
			contentType: "image/jpeg",
		},
	} {
		t.Run("put and get "+scenario.name, func(t *testing.T) {
			id := client.CreateLayout(Layout{}).Id
			defer func() { client.DeleteLayout(id) }()

			client.PutLayer(id, scenario.name, bytes.NewBuffer(scenario.layer))
			gotReadCloser, gotContentType := client.GetLayer(id, scenario.name)
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

	t.Run("put heightmap updates heightmap URL", func(t *testing.T) {
		id := client.CreateLayout(Layout{}).Id
		defer func() { client.DeleteLayout(id) }()
		client.PutLayer(id, "heightmap.jpg", nil)

		layout := client.GetLayout(id)

		got := layout.HeightmapURL
		want := fmt.Sprintf("%s/v1/layouts/%s/heightmap.jpg", baseURL, id)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
