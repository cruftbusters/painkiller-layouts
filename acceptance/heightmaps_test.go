package acceptance

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/maps"
	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestHeightmaps(t *testing.T) {
	listener, port := RandomPortListener()
	go func() {
		http.Serve(listener, maps.Handler())
	}()

	baseURL := fmt.Sprintf("http://localhost:%d", port)
	client := NewClient(t, baseURL)

	t.Run("get missing heightmap", func(t *testing.T) {
		client.GetExpectNotFound("deadbeef")
	})

	t.Run("create and get heightmap", func(t *testing.T) {
		got := client.Create(Metadata{})
		AssertMetadata(t, got, Metadata{Id: got.Id})
		AssertMetadata(t, client.Get(got.Id), got)

		client.Delete(got.Id)
	})

	t.Run("create and get all heightmaps", func(t *testing.T) {
		first := client.Create(Metadata{})
		second := client.Create(Metadata{})

		got := client.GetAll()
		want := []Metadata{first, second}
		AssertAllMetadataUnordered(t, got, want)

		client.Delete(first.Id)
		client.Delete(second.Id)
	})

	t.Run("patch url onto metadata", func(t *testing.T) {
		oldSize, newImageURL := Size{1, 2}, "new image url"
		metadata := client.Create(Metadata{Size: oldSize})

		got := client.Patch(metadata.Id, Metadata{ImageURL: newImageURL})
		want := Metadata{Id: metadata.Id, Size: oldSize, ImageURL: newImageURL}
		AssertMetadata(t, got, want)

		got = client.Get(metadata.Id)
		AssertMetadata(t, got, want)
	})

	t.Run("delete heightmap", func(t *testing.T) {
		metadata := client.Create(Metadata{})
		client.Delete(metadata.Id)
		client.GetExpectNotFound(metadata.Id)
	})
}
