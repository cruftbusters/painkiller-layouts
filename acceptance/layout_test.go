package acceptance

import (
	"fmt"
	"testing"

	"github.com/cruftbusters/painkiller-layouts/layouts"
	. "github.com/cruftbusters/painkiller-layouts/testing"
	. "github.com/cruftbusters/painkiller-layouts/types"
)

func TestLayout(t *testing.T) {
	client, _ := NewTestClient(layouts.Handler)

	t.Run("get missing layout", func(t *testing.T) {
		client.GetLayoutExpectNotFound(t, "deadbeef")
	})

	t.Run("patch missing layout", func(t *testing.T) {
		client.PatchLayoutExpectNotFound(t, "garbotron")
	})

	t.Run("create and get layout", func(t *testing.T) {
		got := client.CreateLayout(t, Layout{})
		defer func() { client.DeleteLayout(t, got.Id) }()
		AssertLayout(t, got, Layout{Id: got.Id})
		AssertLayout(t, client.GetLayout(t, got.Id), got)
	})

	t.Run("delete layout", func(t *testing.T) {
		layout := client.CreateLayout(t, Layout{})
		client.DeleteLayout(t, layout.Id)
		client.GetLayoutExpectNotFound(t, layout.Id)
	})

	t.Run("get all layouts", func(t *testing.T) {
		withHeightmap := client.CreateLayout(t, Layout{HeightmapURL: "heightmap url"})
		defer func() { client.DeleteLayout(t, withHeightmap.Id) }()
		withHillshade := client.CreateLayout(t, Layout{HillshadeURL: "hillshade url"})
		defer func() { client.DeleteLayout(t, withHillshade.Id) }()
		withEverythingElse := client.CreateLayout(t, Layout{
			Size:   Size{Width: 1, Height: 2},
			Bounds: Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
		})
		defer func() { client.DeleteLayout(t, withEverythingElse.Id) }()

		AssertLayoutsUnordered(t,
			client.GetLayouts(t),
			[]Layout{withEverythingElse, withHillshade, withHeightmap},
		)

		t.Run("with no heightmaps", func(t *testing.T) {
			AssertLayoutsUnordered(t,
				client.GetLayoutsWithoutHeightmap(t),
				[]Layout{withEverythingElse, withHillshade},
			)
		})

		t.Run("with heightmap with no hillshade", func(t *testing.T) {
			AssertLayoutsUnordered(t,
				client.GetLayoutsWithHeightmapWithoutHillshade(t),
				[]Layout{withHeightmap},
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
			layout := client.CreateLayout(t,
				Layout{
					Size:         Size{Width: 1, Height: 2},
					Bounds:       Bounds{Left: 3, Top: 4, Right: 5, Bottom: 6},
					HeightmapURL: "old heightmap url",
					HillshadeURL: "old hillshade url",
				},
			)
			defer func() { client.DeleteLayout(t, layout.Id) }()

			got := client.PatchLayout(t, layout.Id, scenario.patch)
			scenario.want(&layout)
			AssertLayout(t, got, layout)

			got = client.GetLayout(t, layout.Id)
			AssertLayout(t, got, layout)
		})
	}
}
