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

	t.Run("put heightmap on non-existant map is not found", func(t *testing.T) {
		client.PutHeightmapExpectNotFound("deadbeef")
	})

	t.Run("put heightmap", func(t *testing.T) {
		id := client.Create(Metadata{}).Id
		client.PutHeightmap(id)
	})
}
