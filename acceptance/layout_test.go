package acceptance

import (
	"testing"

	"github.com/cruftbusters/painkiller-gallery/layouts"
	. "github.com/cruftbusters/painkiller-gallery/testing"
	. "github.com/cruftbusters/painkiller-gallery/types"
)

func TestLayout(t *testing.T) {
	client, _ := NewTestClient(t, layouts.Handler)

	t.Run("get missing map", func(t *testing.T) {
		client.GetExpectNotFound("deadbeef")
	})

	t.Run("create and get map", func(t *testing.T) {
		got := client.Create(Layout{})
		defer func() { client.Delete(got.Id) }()
		AssertLayout(t, got, Layout{Id: got.Id})
		AssertLayout(t, client.Get(got.Id), got)

	})

	t.Run("create and get all maps", func(t *testing.T) {
		first := client.Create(Layout{})
		defer func() { client.Delete(first.Id) }()
		second := client.Create(Layout{})
		defer func() { client.Delete(second.Id) }()

		got := client.GetAll()
		want := []Layout{first, second}
		AssertLayoutsUnordered(t, got, want)
	})

	t.Run("patch missing map", func(t *testing.T) {
		client.PatchExpectNotFound("garbotron")
	})

	t.Run("patch heightmap url onto map", func(t *testing.T) {
		oldSize, newHeightmapURL := Size{Width: 1, Height: 2}, "new heightmap url"
		layout := client.Create(Layout{Size: oldSize})
		defer func() { client.Delete(layout.Id) }()

		got := client.Patch(layout.Id, Layout{HeightmapURL: newHeightmapURL})
		want := Layout{Id: layout.Id, Size: oldSize, HeightmapURL: newHeightmapURL}
		AssertLayout(t, got, want)

		got = client.Get(layout.Id)
		AssertLayout(t, got, want)

	})

	t.Run("filter for maps with no heightmap", func(t *testing.T) {
		withoutHeightmap := client.Create(Layout{})
		defer func() { client.Delete(withoutHeightmap.Id) }()
		withHeightmap := client.Create(Layout{HeightmapURL: "heightmap url"})
		defer func() { client.Delete(withHeightmap.Id) }()

		AssertLayouts(t,
			client.GetAllWithoutHeightmap(),
			[]Layout{withoutHeightmap},
		)
	})

	t.Run("delete map", func(t *testing.T) {
		layout := client.Create(Layout{})
		client.Delete(layout.Id)
		client.GetExpectNotFound(layout.Id)
	})
}
