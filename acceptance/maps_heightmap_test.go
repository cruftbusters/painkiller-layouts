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
	go func() {
		http.Serve(listener, maps.Handler())
	}()

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	client := NewClientV2(t, baseURL)

	t.Run("put heightmap on missing map is not found", func(t *testing.T) {
		client.PutHeightmapExpectNotFound("deadbeef")
	})

	t.Run("get heightmap on missing map is not found", func(t *testing.T) {
		client.GetHeightmapExpectNotFound("deadbeef")
	})

	t.Run("get heightmap is not found", func(t *testing.T) {
		id := client.Create(Metadata{}).Id
		client.GetHeightmapExpectNotFound(id)
	})

	t.Run("put and get heightmap", func(t *testing.T) {
		id, heightmap, contentType := client.Create(Metadata{}).Id, []byte{65, 66, 67}, "image/jpeg"
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
}
