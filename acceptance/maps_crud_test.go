package acceptance

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cruftbusters/painkiller-gallery/maps"
	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestMapsCrud(t *testing.T) {
	listener, port := RandomPortListener()
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	go func() {
		http.Serve(listener, maps.Handler(baseURL))
	}()

	client := NewClientV2(t, baseURL)

	t.Run("get missing map", func(t *testing.T) {
		client.GetExpectNotFound("deadbeef")
	})

	t.Run("create and get map", func(t *testing.T) {
		got := client.Create(Metadata{})
		AssertMetadata(t, got, Metadata{Id: got.Id})
		AssertMetadata(t, client.Get(got.Id), got)

		client.Delete(got.Id)
	})

	t.Run("create and get all maps", func(t *testing.T) {
		first := client.Create(Metadata{})
		second := client.Create(Metadata{})

		got := client.GetAll()
		want := []Metadata{first, second}
		AssertAllMetadataUnordered(t, got, want)

		client.Delete(first.Id)
		client.Delete(second.Id)
	})

	t.Run("patch missing map", func(t *testing.T) {
		client.PatchExpectNotFound("garbotron")
	})

	t.Run("patch heightmap url onto map", func(t *testing.T) {
		oldSize, newHeightmapURL := Size{Width: 1, Height: 2}, "new heightmap url"
		metadata := client.Create(Metadata{Size: oldSize})

		got := client.Patch(metadata.Id, Metadata{HeightmapURL: newHeightmapURL})
		want := Metadata{Id: metadata.Id, Size: oldSize, HeightmapURL: newHeightmapURL}
		AssertMetadata(t, got, want)

		got = client.Get(metadata.Id)
		AssertMetadata(t, got, want)

		client.Delete(metadata.Id)
	})

	t.Run("filter for maps with no heightmap", func(t *testing.T) {
		withoutHeightmap := client.Create(Metadata{})
		withHeightmap := client.Create(Metadata{HeightmapURL: "heightmap url"})

		AssertAllMetadata(t,
			client.GetAllWithoutHeightmap(),
			[]Metadata{withoutHeightmap},
		)

		client.Delete(withoutHeightmap.Id)
		client.Delete(withHeightmap.Id)
	})

	t.Run("delete map", func(t *testing.T) {
		metadata := client.Create(Metadata{})
		client.Delete(metadata.Id)
		client.GetExpectNotFound(metadata.Id)
	})
}
