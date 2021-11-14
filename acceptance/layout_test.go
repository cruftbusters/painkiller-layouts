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
		client.GetLayoutExpectNotFound("deadbeef")
	})

	t.Run("create and get map", func(t *testing.T) {
		got := client.CreateLayout(Layout{})
		defer func() { client.DeleteLayout(got.Id) }()
		AssertLayout(t, got, Layout{Id: got.Id})
		AssertLayout(t, client.GetLayout(got.Id), got)
	})

	t.Run("create and get all maps", func(t *testing.T) {
		first := client.CreateLayout(Layout{})
		defer func() { client.DeleteLayout(first.Id) }()
		second := client.CreateLayout(Layout{})
		defer func() { client.DeleteLayout(second.Id) }()

		got := client.GetLayouts()
		want := []Layout{first, second}
		AssertLayoutsUnordered(t, got, want)
	})

	t.Run("patch missing map", func(t *testing.T) {
		client.PatchLayoutExpectNotFound("garbotron")
	})

	t.Run("patch heightmap url onto map", func(t *testing.T) {
		oldSize, newHeightmapURL := Size{Width: 1, Height: 2}, "new heightmap url"
		layout := client.CreateLayout(Layout{Size: oldSize})
		defer func() { client.DeleteLayout(layout.Id) }()

		got := client.PatchLayout(layout.Id, Layout{HeightmapURL: newHeightmapURL})
		want := Layout{Id: layout.Id, Size: oldSize, HeightmapURL: newHeightmapURL}
		AssertLayout(t, got, want)

		got = client.GetLayout(layout.Id)
		AssertLayout(t, got, want)
	})

	t.Run("filter for maps with no heightmap", func(t *testing.T) {
		withoutHeightmap := client.CreateLayout(Layout{})
		defer func() { client.DeleteLayout(withoutHeightmap.Id) }()
		withHeightmap := client.CreateLayout(Layout{HeightmapURL: "heightmap url"})
		defer func() { client.DeleteLayout(withHeightmap.Id) }()

		AssertLayouts(t,
			client.GetLayoutsWithoutHeightmap(),
			[]Layout{withoutHeightmap},
		)
	})

	t.Run("delete map", func(t *testing.T) {
		layout := client.CreateLayout(Layout{})
		client.DeleteLayout(layout.Id)
		client.GetLayoutExpectNotFound(layout.Id)
	})
}
