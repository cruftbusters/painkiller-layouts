package acceptance

import (
	"fmt"
	"net/http"
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
		id := client.Create(Metadata{}).Id
		client.PutHeightmap(id)
		got := client.GetHeightmap(id)
		want := "fixed bytes"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
