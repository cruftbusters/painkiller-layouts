package acceptance

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/maps"
	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestHeightmap(t *testing.T) {
	listener, port := RandomPortListener()
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	go func() {
		http.Serve(listener, maps.Handler(baseURL))
	}()

	client := NewClientV2(t, baseURL)

	t.Run("put heightmap on missing map is not found", func(t *testing.T) {
		client.PutHeightmapExpectNotFound("deadbeef")
	})

	t.Run("get heightmap on missing map is not found", func(t *testing.T) {
		client.GetHeightmapExpectNotFound("deadbeef")
	})

	t.Run("get heightmap is not found", func(t *testing.T) {
		id := client.Create(Metadata{}).Id
		defer func() { client.Delete(id) }()
		client.GetHeightmapExpectNotFound(id)
	})

	t.Run("put and get heightmap", func(t *testing.T) {
		id, heightmap, contentType := client.Create(Metadata{}).Id, []byte{65, 66, 67}, "image/jpeg"
		defer func() { client.Delete(id) }()
		client.PutHeightmap(id, bytes.NewBuffer(heightmap))
		gotReadCloser, gotContentType := client.GetHeightmap(id)
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
		id := client.Create(Metadata{}).Id
		defer func() { client.Delete(id) }()
		client.PutHeightmap(id, nil)

		metadata := client.Get(id)

		got := metadata.HeightmapURL
		want := fmt.Sprintf("%s/v1/maps/%s/heightmap.jpg", baseURL, id)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
