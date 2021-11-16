package acceptance

import (
	"fmt"
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func TestLayout(t *testing.T) {
	client, _ := NewTestClient(t, layouts.Handler)

	t.Run("get missing layout", func(t *testing.T) {
		client.GetLayoutExpectNotFound("deadbeef")
	})

	t.Run("patch missing layout", func(t *testing.T) {
		client.PatchLayoutExpectNotFound("garbotron")
	})

	t.Run("create and get layout", func(t *testing.T) {
		got := client.CreateLayout(Layout{})
		defer func() { client.DeleteLayout(got.Id) }()
		AssertLayout(t, got, Layout{Id: got.Id})
		AssertLayout(t, client.GetLayout(got.Id), got)
	})

	t.Run("delete layout", func(t *testing.T) {
		layout := client.CreateLayout(Layout{})
		client.DeleteLayout(layout.Id)
		client.GetLayoutExpectNotFound(layout.Id)
	})

	t.Run("get all layouts", func(t *testing.T) {
		withHeightmap := client.CreateLayout(Layout{HeightmapURL: "heightmap url"})
		defer func() { client.DeleteLayout(withHeightmap.Id) }()
		withHillshade := client.CreateLayout(Layout{HillshadeURL: "hillshade url"})
		defer func() { client.DeleteLayout(withHillshade.Id) }()
		withEverythingElse := client.CreateLayout(Layout{
			Size:   Size{Width: 1, Height: 2},
			Bounds: Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
		})
		defer func() { client.DeleteLayout(withEverythingElse.Id) }()

		AssertLayoutsUnordered(t,
			client.GetLayouts(),
			[]Layout{withEverythingElse, withHillshade, withHeightmap},
		)

		t.Run("with no heightmaps", func(t *testing.T) {
			AssertLayoutsUnordered(t,
				client.GetLayoutsWithoutHeightmap(),
				[]Layout{withEverythingElse, withHillshade},
			)
		})

		t.Run("with no hillshades", func(t *testing.T) {
			AssertLayoutsUnordered(t,
				client.GetLayoutsWithoutHillshade(),
				[]Layout{withEverythingElse, withHeightmap},
			)
		})
	})

	for _, scenario := range []struct {
		patch Layout
		want  func(*Layout)
	}{
		{
			patch: Layout{HeightmapURL: "new heightmap url"},
			want:  func(initial *Layout) { initial.HeightmapURL = "new heightmap url" },
		},
		{
			patch: Layout{HillshadeURL: "new hillshade url"},
			want:  func(initial *Layout) { initial.HillshadeURL = "new hillshade url" },
		},
	} {
		t.Run(fmt.Sprintf("patch layout with %+v", scenario.patch), func(t *testing.T) {
			layout := client.CreateLayout(
				Layout{
					Size:         Size{Width: 1, Height: 2},
					Bounds:       Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
					HeightmapURL: "old heightmap url",
					HillshadeURL: "old hillshade url",
				},
			)
			defer func() { client.DeleteLayout(layout.Id) }()

			got := client.PatchLayout(layout.Id, scenario.patch)
			scenario.want(&layout)
			AssertLayout(t, got, layout)

			got = client.GetLayout(layout.Id)
			AssertLayout(t, got, layout)
		})
	}
}
