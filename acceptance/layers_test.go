package acceptance

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func TestLayers(t *testing.T) {
	client, baseURL := NewTestClient(t, layouts.Handler)

	t.Run("put heightmap on missing map is not found", func(t *testing.T) {
		client.PutLayerExpectNotFound("deadbeef", "heightmap.jpg")
	})

	t.Run("get heightmap on missing map is not found", func(t *testing.T) {
		client.GetLayerExpectNotFound("deadbeef", "heightmap.jpg")
	})

	t.Run("get heightmap is not found", func(t *testing.T) {
		id := client.CreateLayout(Layout{}).Id
		defer func() { client.DeleteLayout(id) }()
		client.GetLayerExpectNotFound(id, "heightmap.jpg")
	})

	t.Run("put and get heightmap", func(t *testing.T) {
		id, heightmap, contentType := client.CreateLayout(Layout{}).Id, []byte{65, 66, 67}, "image/jpeg"
		defer func() { client.DeleteLayout(id) }()
		client.PutLayer(id, "heightmap.jpg", bytes.NewBuffer(heightmap))
		gotReadCloser, gotContentType := client.GetLayer(id, "heightmap.jpg")
		got, err := io.ReadAll(gotReadCloser)
		AssertNoError(t, err)
		want := heightmap
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
		if gotContentType != contentType {
			t.Errorf("got %s want %s", gotContentType, contentType)
		}
	})

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
