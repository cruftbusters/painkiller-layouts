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

func TestLayoutHeightmap(t *testing.T) {
	client, baseURL := NewTestClient(t, layouts.Handler)

	t.Run("put heightmap on missing map is not found", func(t *testing.T) {
		client.PutLayoutHeightmapExpectNotFound("deadbeef")
	})

	t.Run("get heightmap on missing map is not found", func(t *testing.T) {
		client.GetLayoutHeightmapExpectNotFound("deadbeef")
	})

	t.Run("get heightmap is not found", func(t *testing.T) {
		id := client.CreateLayout(Layout{}).Id
		defer func() { client.DeleteLayout(id) }()
		client.GetLayoutHeightmapExpectNotFound(id)
	})

	t.Run("put and get heightmap", func(t *testing.T) {
		id, heightmap, contentType := client.CreateLayout(Layout{}).Id, []byte{65, 66, 67}, "image/jpeg"
		defer func() { client.DeleteLayout(id) }()
		client.PutLayoutHeightmap(id, bytes.NewBuffer(heightmap))
		gotReadCloser, gotContentType := client.GetLayoutHeightmap(id)
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
		client.PutLayoutHeightmap(id, nil)

		layout := client.Get(id)

		got := layout.HeightmapURL
		want := fmt.Sprintf("%s/v1/maps/%s/heightmap.jpg", baseURL, id)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
